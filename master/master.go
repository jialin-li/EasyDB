package main

import (
	//"EasyDB/client"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const clientPath = "/u/wlsaidhi/go/bin/client"

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')

		client := exec.Command(clientPath)
		client.Start()

		switch strs := strings.Split(text, " "); strs[0] {
		case "joinServer":
			fmt.Println(strs[1])
		case "killServer":
			fmt.Println(strs[1])
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
