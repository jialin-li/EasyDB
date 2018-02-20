package main

import (
	"bufio"
	// "errors"
	"flag"
	"fmt"
	"github.com/jialin-li/EasyDB/shared"
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

var serverIds map[int]int
var clientIds map[int]int
var clientId = 0
var serverId = 0

const serverPath = "./server"
const clientPath = "./client"

var clientCalls map[int]*rpcClient
var serverCalls map[int]*rpcClient

var wg sync.WaitGroup

var term bool

func main() {
	termPtr := flag.Bool("term", false, "run as program")

	flag.Parse()
	if *termPtr {
		term = true
	}

	serverIds = make(map[int]int)
	clientIds = make(map[int]int)

	clientCalls = make(map[int]*rpcClient)
	serverCalls = make(map[int]*rpcClient)

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
			var unused string
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
			fmt.Println(strs[1])
		case "joinClient":
			var unused string
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
			var unused string
			var id1, id2 int
			_, err := fmt.Sscanf(text, "%s %d %d", &unused, &id1, &id2)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println(id1, id2)
		case "createConnection":
			var unused string
			var id1, id2 int
			_, err := fmt.Sscanf(text, "%s %d %d", &unused, &id1, &id2)
			if err != nil {
				log.Println(err)
				continue
			}
			createConnection(id1, id2)

		case "stabilize":
			fmt.Println(strs[1])
		case "printStore":
			var unused string
			var serverId int
			_, err := fmt.Sscanf(
				text, "%s %d", &unused, &serverId)
			if err != nil {
				log.Println(err)
				continue
			}
			// translate to internal server id
			serverId = clientIds[serverId]
			printStore(serverId)

		case "put":
			var unused, key, value string
			var clientId int
			_, err := fmt.Sscanf(
				text, "%s %d %s %s", &unused, &clientId, &key, &value)
			if err != nil {
				log.Println(err)
				continue
			}
			// translate to internal client id
			clientId = clientIds[clientId]
			put(clientId, key, value)

		case "get":
			var unused, key string
			var clientId int
			_, err := fmt.Sscanf(
				text, "%s %d %s", &unused, &clientId, &key)
			if err != nil {
				log.Println(err)
				continue
			}
			// translate to internal client id
			clientId = clientIds[clientId]
			get(clientId, key)

		default:
			fmt.Println("bad command")
		}
	}
}

func joinServer(id int) error {
	// check if a server with id already exists
	if _, ok := serverIds[id]; ok {
		return fmt.Errorf("joinServer: server id %d already exists!", id)
	}

	// update to our mapping
	sid := getServerId()
	serverIds[id] = sid

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
	if _, ok := clientIds[clientId]; ok {
		return fmt.Errorf("joinClient: client id %d already exists!", clientId)
	}

	// get the server id
	if _, ok := serverIds[serverId]; !ok {
		return fmt.Errorf("joinClient: server id %d does not exist!", serverId)
	}

	// update to our mapping
	cid := getClientId()
	clientIds[clientId] = cid

	// start a new client
	client := exec.Command(clientPath, strconv.Itoa(cid), strconv.Itoa(serverIds[serverId]))
	err := client.Start()

	// wait for client to notify us before proceeding
	if err == nil {
		wg.Add(1)
		// TODO: add a timeout?
		wg.Wait()
	}
	return err
}

func createConnection(id1, id2 int) {
	//rpcCalls[clientId].connect("connecting", "test", serverId)
}

func breakConnection(id1, id2 int) {
}

func printStore(serverId int) {
	value := serverCalls[serverId].printStore()
	// parse value
	store := strings.Split(value, " ")
	for _, v := range store {
		fmt.Println(v)
	}
}
func put(clientId int, key, value string) error {
	clientCalls[clientId].put(key, value)
	return nil
}

func get(clientId int, key string) error {
	value := clientCalls[clientId].get(key)
	fmt.Printf("Retrieved: %v:%v\n", key, value)
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

func getIdType(id int) int {
	if _, ok := serverIds[id]; ok {
		return shared.ServerType
	} else if _, ok := clientIds[id]; ok {
		return shared.ClientType
	}
	return -1
}
