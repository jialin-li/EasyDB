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

	"github.com/jialin-li/EasyDB/shared" //Path to the package contains shared struct
)

var wg sync.WaitGroup

//type dbValue struct {
//value string
//time  shared.Time
//}

var db map[string]*shared.DbValue
var serverListener net.Listener
var serverCalls map[int]*rpcClient
var serverId = -1
var term bool

func main() {
	//f, err := os.OpenFile("output", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//defer f.Close()
	//log.SetOutput(f)

	// Use the -term flag to run  the server as a command line program. Server
	// will wait for commands from stdin. Useful for debugging and for a real
	// distributed system.
	termPtr := flag.Bool("term", false, "run as program")

	args := os.Args[1:]

	flag.Parse()
	if *termPtr {
		term = true
		args = os.Args[2:]
	}

	var err error
	serverId, err = strconv.Atoi(args[0])
	if err != nil {
		log.Println(err)
	}

	// set up key value store
	db = make(map[string]*shared.DbValue)

	serverCalls = make(map[int]*rpcClient)

	listen(shared.BasePort + serverId)

	// Tries to connect to localhost:1234 (The port on which master's rpc
	// server is listening)
	conn, _ := shared.Dial(shared.MasterPort)

	masterCall := &rpcClient{client: rpc.NewClient(conn)}
	masterCall.notify(serverId)

	if term {
		parseCommands()
	}

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
			printStore()
		case "put":
			fmt.Println(strs[1])
		case "get":
			fmt.Println(strs[1])
		default:
			fmt.Println(strs[0])
			fmt.Println("bad command")
		}
	}
}

func printStore() {
	for k, v := range db {
		fmt.Println(k + ":" + v.Value)
	}
}

func dumpStore() string {
	if len(db) == 0 {
		return ""
	}

	var store string
	for k, v := range db {
		store += k + ":" + v.Value + " "
	}
	// cut off last space
	return store[:len(store)-1]
}

func listen(port int) error {
	// create an instance of struct that implements shared.Server interface
	serverInterface := new(KVServer)

	// register a new rpc server
	rpcServer := rpc.NewServer()
	rpcServer.Register(serverInterface)

	// Listen for incoming tcp packets on specified port.
	serverListener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	// if everything goes well, we can call Accept in a go routine and add to
	// waitgroup to allow main to block after setup
	if err == nil {
		wg.Add(1)
		// TODO: handle failure more gracefully
		go func() {
			defer wg.Done()
			rpcServer.Accept(serverListener)
		}()
	}
	return err
}
