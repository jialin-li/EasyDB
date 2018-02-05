package main

import (
	"fmt"
	"github.com/jialin-li/EasyDB/shared"
)

type Server int

func (t *Server) Terminate() error {
	return nil
}
func (t *Server) Connect(args *shared.Args, reply *shared.Response) error {
	*reply = shared.Response{"it worked"}
	fmt.Println(args.Msg, args.Key, args.Value)
	return nil
}

// Terminate a server
// Ask the server to connect to another server

// Send back the current key value store
func (t *Server) DumpStore(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Disconnect if we are clients of any other servers
func (t *Server) Disconnect(args *shared.Args, reply *shared.Response) error {
	return nil
}

// No more new write requests will be send until the next write
func (t *Server) Stabilize(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Put a KV pair
func (t *Server) Put(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Get a Value based on a key
func (t *Server) Get(args *shared.Args, reply *shared.Response) error {
	return nil
}
