package main

import (
	//"EasyDB/client"
	"bufio"
	"errors"
	"fmt"
	"github.com/jialin-li/EasyDB/shared"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Connection struct {
	Type int
	Port string
	Conn net.Listener
}

var connections map[int]Connection

const serverPath = "./server"
const clientPath = "./client"

var port = 1234

type Server int

func (t *Server) Connect(args *shared.Args, reply *shared.Response) error {
	*reply = shared.Response{"it worked"}
	fmt.Println(args.Msg, args.Key, args.Value)
	// add client connection to map
	//connections[clientId] = Connection{1, port, conn}
	return nil
}

func main() {
	connections = make(map[int]Connection)

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
			fmt.Println(strs[1])
		case "get":
			fmt.Println(strs[1])
		default:
			fmt.Println("bad command")
		}
	}
}

func registerServer(server *rpc.Server, s shared.Server) {
	server.Register(s)
}

func joinServer(id int) error {

	// check if a server with id already exists
	if _, ok := connections[id]; !ok {
		return errors.New("joinServer: server id already exists!")
	}

	// create an instance of struct that implements Server interface
	serverInterface := new(Server)

	// register a new rpc server
	rpcServer := rpc.NewServer()
	registerServer(rpcServer, serverInterface)

	// start a new client
	server := exec.Command(serverPath)
	err := server.Start()

	return err
}

func joinClient(clientId, serverId int) error {

	// check if a client with id already exists
	if _, ok := connections[clientId]; ok {
		return errors.New("joinClient: client id already exists!")
	}

	// start a new client
	client := exec.Command(clientPath)
	err := client.Start()
	if err != nil {
		return err
	}

	return err
}

func listen(port int) error {
	// create an instance of struct that implements Server interface
	serverInterface := new(Server)

	// register a new rpc server
	rpcServer := rpc.NewServer()
	registerServer(rpcServer, serverInterface)

	// Listen for incoming tcp packets on specified port.
	conn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		go rpcServer.Accept(conn)
	}
	return err
}
