package main

import (
	//"EasyDB/client"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"
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

	// Tries to connect to localhost:1234 (The port on which rpc server is
	// listening)
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Connection:", err)
	} else {
		fmt.Println("sucess?")
	}

	client := &Client{client: rpc.NewClient(conn)}
	fmt.Println(client.Connect("please work"))
}

func parseCommands() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')

		switch strs := strings.Split(text, " "); strs[0] {
		case "joinClient":
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

func client(id int) int {
	for i := 0; i < id; i++ {
		fmt.Println("in client")
	}
	return 0

}
