package main

import (
	"fmt"

	"github.com/jialin-li/EasyDB/shared"
)

type KVServer struct {
}

// Terminate a server
func (*KVServer) Terminate() error {
	return nil
}

// Ask the server to connect to another server
func (*KVServer) Connect(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Send back the current key value store
func (*KVServer) DumpStore(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Disconnect if we are clients of any other servers
func (*KVServer) Disconnect(args *shared.Args, reply *shared.Response) error {
	return nil
}

// No more new write requests will be send until the next write
func (*KVServer) Stabilize(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Put a KV pair
func (*KVServer) Put(args *shared.Args, reply *shared.Response) error {
	fmt.Printf("%s:%s \n", args.Key, args.Value)
	reply.Result = "hello friend"
	return nil
}

// Get a Value based on a key
func (*KVServer) Get(args *shared.Args, reply *shared.Response) error {
	return nil
}
