package main

import (
	"bufio"
	"flag"
	"fmt"
	//	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	// only added for testing
	"github.com/jialin-li/EasyDB/shared"
)

func main() {
	// Use the -term flag to run  the client as a command line program. Client
	// will wait for commands from stdin. Useful for debugging and for a real
	// distributed system.
	termPtr := flag.Bool("term", false, "run as program")

	flag.Parse()
	if *termPtr {
		parseCommands()
	}

	// Tries to connect to localhost:1234 (The port on which rpc server is listening)
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf("Connectiong: %e \n", err)
	}

	// Create a struct, that mimics all methods provided by interface.
	// It is not compulsory, we are doing it here, just to simulate a traditional method call.
	client := &KVClient{client: rpc.NewClient(conn)}
	client.Notify("Notifying master")
	test(client)
}

func test(c *KVClient) {
	// testing client
	args := &shared.Args{"", "Hi", "World!"}
	resp := &shared.Response{}
	if err := c.Put(args, resp); err != nil {
		fmt.Printf("Error returned from client.Put %v \n", err)
	} else {
		fmt.Printf("Heard back from server %s \n", resp.Result)
	}
}

func parseCommands() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')

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
