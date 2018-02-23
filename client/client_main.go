package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"sync"

	// only added for testing
	"github.com/jialin-li/EasyDB/shared"
)

var serverCalls map[int]*rpcClient
var keyTimes map[string]*shared.Time

// tracking our own id
var clientId int

var wg sync.WaitGroup

func main() {
	//f, err := os.OpenFile("output", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//defer f.Close()
	//log.SetOutput(f)

	// Use the -term flag to run  the client as a command line program. Client
	// will wait for commands from stdin. Useful for debugging and for a real
	// distributed system.
	termPtr := flag.Bool("term", false, "run as program")

	flag.Parse()
	if *termPtr {
		parseCommands()
	}
	args := os.Args[1:]

	// var clientId, serverId int

	// if we fail to extract our id, we should probably exit instead
	clientId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Println(err)
	}

	serverCalls = make(map[int]*rpcClient)
	keyTimes = make(map[string]*shared.Time)

	serverId, err := strconv.Atoi(args[1])
	if err != nil {
		log.Println(err)
	}
	setupConn(serverId)

	// Listen for connections from master and servers
	listen(shared.BasePort + clientId)

	// Tries to connect to localhost:1234 (The port on which master's rpc
	// server is listening
	conn, _ := shared.Dial(shared.MasterPort)

	// Create a struct, that mimics all methods provided by interface.
	// It is not compulsory, we are doing it here, just to simulate a traditional method call.
	masterCall := &rpcClient{client: rpc.NewClient(conn)}
	masterCall.notify(clientId)

	// now we block
	wg.Wait()
}

func parseCommands() {
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
		case "breakConnection":
			fmt.Println(strs[1])
		case "createConnection":
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

func listen(port int) error {
	// create an instance of struct that implements shared.Client interface
	serverInterface := new(KVClient)

	// register a new rpc server
	rpcServer := rpc.NewServer()
	rpcServer.Register(serverInterface)

	// Listen for incoming tcp packets on specified port.
	conn, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	// if everything goes well, we can call Accept in a go routine and add to
	// waitgroup to allow main to block after setup
	if err == nil {
		wg.Add(1)
		// TODO: handle failure more gracefully
		go func() {
			defer wg.Done()
			rpcServer.Accept(conn)
		}()
	}
	return err
}
