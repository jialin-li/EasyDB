package main

import (
	//"EasyDB/client"
	"bufio"
	// "flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"strings"

	"github.com/jialin-li/EasyDB/shared"
)

type KVClient struct {
	client *rpc.Client
}

func main() {

	// Use the -term flag to run  the client as a command line program. Client
	// will wait for commands from stdin. Useful for debugging and for a real
	// distributed system.
	// termPtr := flag.Bool("term", false, "run as program")

	// flag.Parse()
	// if *termPtr {
	// 	parseCommands()
	// }
	fmt.Println("In main")

	// Tries to connect to localhost:1234 (The port on which rpc server is listening)
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Printf("Connectiong: %e \n", err)
	}

	// Create a struct, that mimics all methods provided by interface.
	// It is not compulsory, we are doing it here, just to simulate a traditional method call.
	client := &KVClient{client: rpc.NewClient(conn)}

	fmt.Println(client.Put("Hi", "World!"))

	// client(5)
}

func (t *KVClient) Put(a, b string) shared.Response {
	args := &shared.Args{"", a, b}
	var reply shared.Response
	err := t.client.Call("KVServer.Put", args, &reply)
	if err != nil {
		fmt.Printf("kvserver error: %e \n", err)
	}
	return reply
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
