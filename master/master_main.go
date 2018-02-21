package main

import (
	"bufio"
	// "errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type connection struct {
	Type int
	Port int
	Conn net.Conn
}

var IdMap map[int]int

// var clientIds map[int]int
var clientId = 0
var serverId = 0

//var clientConnections map[int]connection
// var serverConnections map[int]connection
const serverPath = "./server"
const clientPath = "./client"

var port = 1234

// var clientCalls map[int]*rpcClient
var conns map[int]*rpcClient

var wg sync.WaitGroup

var term bool

func main() {
	//clientConnections = make(map[int]connection)
	// serverConnections = make(map[int]connection)

	termPtr := flag.Bool("term", false, "run as program")

	flag.Parse()
	if *termPtr {
		term = true
	}

	IdMap = make(map[int]int)
	// clientIds = make(map[int]int)

	conns = make(map[int]*rpcClient)
	// serverCalls = make(map[int]*rpcClient)

	listen(port)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')

		if err != nil {
			log.Println(err)
			break
		}

		// remove the newline character
		text = text[:len(text)-1]

		switch strs := strings.Split(text, " "); strs[0] {
		case "joinServer":
			var id int
			var err error
			if id, err = strconv.Atoi(strs[1]); err != nil {
				log.Println(err)
				continue
			}

			if err = joinServer(id); err != nil {
				log.Println(err)
				continue
			}
		case "killServer":
			if id, err := strconv.Atoi(strs[1]); err == nil {
				sid := -1
				if !term {
					sid = IdMap[id]
				}
				killServer(sid, id)
			} else {
				log.Println(err)
			}
		case "joinClient":
			var clientId, serverId int
			var err error
			if clientId, err = strconv.Atoi(strs[1]); err != nil {
				log.Println(err)
				continue
			}

			if serverId, err = strconv.Atoi(strs[2]); err != nil {
				log.Println(err)
				continue
			}

			if err = joinClient(clientId, serverId); err != nil {
				log.Println(err)
				continue
			}
		case "breakConnection":
			fmt.Println(strs[1])
		case "createConnection":
			fmt.Println(strs[1])
			fmt.Println(strs[2])
			var clientId, serverId int
			var err error
			if clientId, err = strconv.Atoi(strs[1]); err != nil {
				log.Println(err)
				continue
			}

			if serverId, err = strconv.Atoi(strs[2]); err != nil {
				log.Println(err)
				continue
			}

			conns[clientId].connect("connecting", "test", serverId)
			//if err = joinClient(clientId, serverId); err != nil {
			//log.Println(err)
			//continue
			//}
		case "stabilize":
			fmt.Println(strs[1])
		case "printStore":
			if serverId, err := strconv.Atoi(strs[1]); err == nil {
				if !term {
					serverId = IdMap[serverId]
				}
				printStore(serverId)
			} else {
				log.Println(err)
			}
		case "put":
			if clientId, err := strconv.Atoi(strs[1]); err == nil {
				if !term {
					clientId = IdMap[clientId]
				}
				put(clientId, strs[2], strs[3])
			} else {
				log.Println(err)
			}
		case "get":
			if clientId, err := strconv.Atoi(strs[1]); err == nil {
				if !term {
					clientId = IdMap[clientId]
				}
				get(clientId, strs[2])
			} else {
				log.Println(err)
			}
		default:
			fmt.Println("bad command")
		}
	}
}

func joinServer(id int) error {
	// check if a server with id already exists
	if _, ok := IdMap[id]; ok {
		return fmt.Errorf("joinServer: server id %d already exists!", id)
	}

	// update to our mapping
	sid := getServerId()
	IdMap[id] = sid

	// start a new server
	server := exec.Command(serverPath, strconv.Itoa(sid))
	err := server.Start()

	// wait for server to notify us before proceeding
	if err == nil {
		// NOTE: underlying assumption is that only one Notify will come in at a time
		// this is fine, just need to keep it in mind when we do things that may require
		// multiple servers to coordinate (ex. stablize)
		wg.Add(1)
		// TODO: add a timeout?
		wg.Wait()
	}
	return err
}

func joinClient(clientId, serverId int) error {
	// check if a client with id already exists
	if _, ok := IdMap[clientId]; ok {
		return fmt.Errorf("joinClient: client id %d already exists!", clientId)
	}

	// get the server id
	if _, ok := IdMap[serverId]; !ok {
		return fmt.Errorf("joinClient: server id %d does not exist!", serverId)
	}

	// update to our mapping
	cid := getClientId()
	IdMap[clientId] = cid

	// start a new client
	client := exec.Command(clientPath, strconv.Itoa(cid), strconv.Itoa(IdMap[serverId]))
	err := client.Start()

	// wait for client to notify us before proceeding
	if err == nil {
		wg.Add(1)
		// TODO: add a timeout?
		wg.Wait()
	}
	return err
}

func getClientId() int {
	clientId++
	return clientId - 1
}

func getServerId() int {
	serverId++
	return serverId - 1
}

func printStore(serverId int) {
	value := conns[serverId].printStore()
	// parse value
	store := strings.Split(value, " ")
	for _, v := range store {
		fmt.Println(v)
	}
}

func put(clientId int, key, value string) error {
	if conn, ok := conns[clientId]; ok {
		return conn.put(key, value)
	}
	return fmt.Errorf("Put: Client id does not exist")
}

func get(clientId int, key string) error {
	if conn, ok := conns[clientId]; ok {
		value := conn.get(key)
		fmt.Printf("%v:%v\n", key, value)
		return nil
	}
	return fmt.Errorf("Get: Client id does not exist")
}

func killServer(serverId, id int) error {
	delete(IdMap, id)
	sid := serverId
	if sid == -1 {
		sid = id
	}

	if conn, ok := conns[sid]; ok {
		delete(conns, sid)
		conn.kill()
	}
	return fmt.Errorf("Get: server id does not exist")
}

func listen(port int) error {
	// create an instance of struct that implements Master interface
	serverInterface := new(Master)

	// register a new rpc server
	rpcServer := rpc.NewServer()
	rpcServer.Register(serverInterface)

	// Listen for incoming tcp packets on specified port.
	conn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		go rpcServer.Accept(conn)
	}
	return err
}
