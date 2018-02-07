package main

import (
	//"errors"
	"bufio"
	// "flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"
	// "github.com/jialin-li/EasyDB/shared" //Path to the package contains shared struct
)

func main() {

	// Use the -term flag to run  the server as a command line program. Server
	// will wait for commands from stdin. Useful for debugging and for a real
	// distributed system.
	// termPtr := flag.Bool("term", false, "run as program")

	// flag.Parse()
	// if *termPtr {
	// 	parseCommands()
	// }
	fmt.Println("In server main")

	//Creating an instance of struct which implement Server interface
	kv := new(KVServer)

	fmt.Println("Created new server")

	// Register a new rpc server (In most cases, you will use default server only)
	// And register struct we created above by name "kv"
	// The wrapper method here ensures that only structs which implement Arith interface
	// are allowed to register themselves.
	server := rpc.NewServer()
	server.RegisterName("KVServer", kv)
	// Get the port from OS args
	// Listen for incoming tcp packets on specified port.
	fmt.Println("Registered new server, about to listen ")

	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		fmt.Printf("listen error: \n", e)
	}

	// This statement links rpc server to the socket, and allows rpc server to accept
	// rpc request coming from that socket.

	server.Accept(l)
}

func parseCommands() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')

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
