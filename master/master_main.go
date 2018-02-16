package main

import (
	"bufio"
	//"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
var clientConnections map[int]connection
var serverConnections map[int]connection

const serverPath = "./server"
const clientPath = "./client"

var port = 1234

var remoteCall *rpcClient

func main() {
	clientConnections = make(map[int]connection)
	serverConnections = make(map[int]connection)

	listen(port)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
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
		case "stabilize":
			fmt.Println(strs[1])
		case "printStore":
			fmt.Println(strs[1])
		case "put":
			remoteCall.put("Put request to client", "key", "new value")
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

	return err
}

func joinClient(clientId, serverId int) error {

	// check if a client with id already exists
	//if _, ok := clientConnections[clientId]; ok {
	//return errors.New("joinClient: client id already exists!")
	//}

	// start a new client
	client := exec.Command(clientPath, strconv.Itoa(getClientId()))
	err := client.Start()

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
