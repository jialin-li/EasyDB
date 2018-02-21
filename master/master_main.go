package main

import (
	"bufio"
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

	"github.com/jialin-li/EasyDB/shared"
)

type connection struct {
	Type int
	Port int
	Conn net.Conn
}

var IdMap map[int]int
var conns map[int]*rpcClient

var clientId = shared.ClientStart
var serverId = shared.ServerStart

const serverPath = "./server"
const clientPath = "./client"

var port = 1234

var wg sync.WaitGroup

var term bool

func main() {
	// used for scanf
	var unused string
	termPtr := flag.Bool("term", false, "run as program")

	flag.Parse()
	if *termPtr {
		term = true
	}

	IdMap = make(map[int]int)
	conns = make(map[int]*rpcClient)

	listen(shared.MasterPort)

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

			_, err := fmt.Sscanf(text, "%s %d", &unused, &id)
			if err != nil {
				log.Println(err)
				continue
			}

			if err = joinServer(id); err != nil {
				log.Println(err)
				continue
			}
		case "killServer":
			var id int

			_, err := fmt.Sscanf(text, "%s %d", &unused, &id)
			if err != nil {
				log.Println(err)
				continue
			}

			if err = killServer(id); err != nil {
				log.Println(err)
				continue
			}

		case "joinClient":
			var clientId, serverId int

			_, err := fmt.Sscanf(
				text, "%s %d %d", &unused, &clientId, &serverId)
			if err != nil {
				log.Println(err)
				continue
			}

			if err = joinClient(clientId, serverId); err != nil {
				log.Println(err)
				continue
			}
		case "breakConnection":
			var id1, id2 int
			_, err := fmt.Sscanf(text, "%s %d %d", &unused, &id1, &id2)
			if err != nil {
				log.Println(err)
				continue
			}
			breakConnection(id1, id2)

		case "createConnection":
			var id1, id2 int
			_, err := fmt.Sscanf(text, "%s %d %d", &unused, &id1, &id2)
			if err != nil {
				log.Println(err)
				continue
			}

			if err = createConnection(id1, id2); err != nil {
				log.Println(err)
			}

		case "stabilize":
			fmt.Println(strs[1])
		case "printStore":
			var serverId int
			_, err := fmt.Sscanf(
				text, "%s %d", &unused, &serverId)
			if err != nil {
				log.Println(err)
				continue
			}
			// translate to internal server id
			serverId = IdMap[serverId]
			printStore(serverId)

		case "put":
			var key, value string
			var clientId int
			_, err := fmt.Sscanf(
				text, "%s %d %s %s", &unused, &clientId, &key, &value)
			if err != nil {
				log.Println(err)
				continue
			}
			// translate to internal client id
			clientId = IdMap[clientId]
			put(clientId, key, value)

		case "get":
			var key string
			var clientId int
			_, err := fmt.Sscanf(
				text, "%s %d %s", &unused, &clientId, &key)
			if err != nil {
				log.Println(err)
				continue
			}
			// translate to internal client id
			clientId = IdMap[clientId]
			get(clientId, key)

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
	client := exec.Command(
		clientPath,
		strconv.Itoa(cid),
		strconv.Itoa(IdMap[serverId]))
	err := client.Start()

	// wait for client to notify us before proceeding
	if err == nil {
		wg.Add(1)
		// TODO: add a timeout?
		wg.Wait()
	}
	return err
}

func createConnection(id1, id2 int) error {
	id1, ok := IdMap[id1]
	if !ok {
		return fmt.Errorf("id1 is not a valid id %d ", id1)
	}

	id2, ok = IdMap[id2]
	if !ok {
		return fmt.Errorf("id2 is not a valid id %d ", id2)
	}

	if isClientId(id1) {
		return conns[id1].clientConnect(id2)
	} else if isClientId(id2) {
		return conns[id2].clientConnect(id1)
	}
	return conns[id1].serverConnect(id2)
}

func breakConnection(id1, id2 int) error {
	id1, ok := IdMap[id1]
	if !ok {
		return fmt.Errorf("id1 is not a valid id %d ", id1)
	}

	id2, ok = IdMap[id2]
	if !ok {
		return fmt.Errorf("id2 is not a valid id %d ", id2)
	}

	if isClientId(id1) {
		// fmt.Println("In here!")
		return conns[id1].clientDisconnect(id2)
	} else if isClientId(id2) {
		return conns[id2].clientDisconnect(id1)
	}
	return conns[id1].serverDisconnect(id2)
}

func printStore(serverId int) {
	value := conns[serverId].printStore()
	// parse value
	if value != "" {
		store := strings.Split(value, " ")
		for _, v := range store {
			fmt.Println(v)
		}
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

func killServer(id int) error {
	sid, ok := IdMap[id]
	if !ok {
		return fmt.Errorf("KillServer: server id does not exist")
	}

	defer delete(IdMap, id)

	if conn, ok := conns[sid]; ok {
		conn.kill()
		delete(conns, sid)
	} else {
		return fmt.Errorf("KillServer: server connection does not exist")
	}

	return nil
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

func getClientId() int {
	clientId++
	return clientId - 1
}

func getServerId() int {
	serverId++
	return serverId - 1
}

func isClientId(id int) bool {
	return id < shared.ServerStart
}
