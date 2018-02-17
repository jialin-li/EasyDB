package main

import (
	"bufio"
	"errors"
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

var serverIds map[int]int
var clientIds map[int]int
var clientId = 0
var serverId = 0

//var clientConnections map[int]connection
var serverConnections map[int]connection

const serverPath = "./server"
const clientPath = "./client"

var port = 1234

var clientCalls map[int]*rpcClient

//var serverCalls map[int]*rpcClient

var wg sync.WaitGroup

func main() {
	//clientConnections = make(map[int]connection)
	serverConnections = make(map[int]connection)

	clientCalls = make(map[int]*rpcClient)

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
			fmt.Println(strs[1])
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

			clientCalls[clientId].connect("connecting", "test", serverId)
			//if err = joinClient(clientId, serverId); err != nil {
			//log.Println(err)
			//continue
			//}
		case "stabilize":
			fmt.Println(strs[1])
		case "printStore":
			fmt.Println(strs[1])
		case "put":
			if clientId, err := strconv.Atoi(strs[1]); err == nil {
				put(clientId, strs[2], strs[3])
			} else {
				log.Println(err)
			}
		case "get":
			fmt.Println(strs[1])
		default:
			fmt.Println("bad command")
		}
	}
}

func joinServer(id int) error {

	// check if a server with id already exists
	//if _, ok := serverConnections[id]; !ok {
	//return errors.New("joinServer: server id already exists!")
	//}

	// start a new client
	server := exec.Command(serverPath)
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
	if _, ok := clientCalls[clientId]; ok {
		return errors.New("joinClient: client id already exists!")
	}

	// start a new client
	client := exec.Command(clientPath, strconv.Itoa(getClientId()))
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

func put(clientId int, key, value string) error {
	clientCalls[clientId].put("Calling put", key, value)
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
