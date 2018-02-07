package main

import (
	//"EasyDB/client"
	"bufio"
	"fmt"
	"github.com/jialin-li/EasyDB/shared"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
)

const clientPath = "./client"

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')

		switch strs := strings.Split(text, " "); strs[0] {
		case "joinServer":
			fmt.Println(strs[1])
		case "killServer":
			fmt.Println(strs[1])
		case "joinClient":
			joinClient()
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

func registerServer(server *rpc.Server, s shared.Server) {
	server.Register(s)
}

func joinClient() {

	server := new(Server)

	s := rpc.NewServer()
	registerServer(s, server)

	// start a new client
	client := exec.Command(clientPath)
	client.Start()

	fmt.Println("listening")
	// Listen for incoming tcp packets on specified port.
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	} else {
		fmt.Println("success?")
	}

	s.Accept(l)
}
