package main

import (
	//"errors"
	"log"
	"net"
	"net/rpc"
	// "github.com/jialin-li/EasyDB/shared" //Path to the package contains shared struct
)

func main() {

	//Creating an instance of struct which implement Arith interface
	kv := new(KVServer)

	// Register a new rpc server (In most cases, you will use default server only)
	// And register struct we created above by name "Arith"
	// The wrapper method here ensures that only structs which implement Arith interface
	// are allowed to register themselves.
	server := rpc.NewServer()
	server.RegisterName("KVServer", kv)
	// Get the port from OS args
	// Listen for incoming tcp packets on specified port.
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	// This statement links rpc server to the socket, and allows rpc server to accept
	// rpc request coming from that socket.
	server.Accept(l)
}
