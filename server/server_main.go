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
	"strconv"
	"strings"
	"sync"

	"github.com/jialin-li/EasyDB/shared" //Path to the package contains shared struct
)

var wg sync.WaitGroup

// temp key value store for testing
type dbValue struct {
	value string
	time  [shared.ClockLen]int
}

var db map[string]*dbValue

func main() {

	// Use the -term flag to run  the server as a command line program. Server
	// will wait for commands from stdin. Useful for debugging and for a real
	// distributed system.
	termPtr := flag.Bool("term", false, "run as program")

	flag.Parse()
	if *termPtr {
		parseCommands()
	}
	args := os.Args[1:]

	serverId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Println(err)
	}

	// set up key value store
	db = make(map[string]dbValue)

	listen(shared.ServerPort + serverId)

	// Tries to connect to localhost:1234 (The port on which master's rpc
	// server is listening)
	conn, _ := shared.Dial(shared.MasterPort)

	masterCall := &rpcClient{client: rpc.NewClient(conn)}
	masterCall.notify(serverId)

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

		switch strs := strings.Split(text, " "); strs[0] {
		case "joinServer":
			fmt.Println(strs[1])
		case "killServer":
			fmt.Println(strs[1])
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

func listen(port int) error {
	// create an instance of struct that implements shared.Server interface
	serverInterface := new(KVServer)

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
