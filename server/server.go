package main

import (
	"fmt"
	"github.com/jialin-li/EasyDB/shared"
	"log"
	"net/rpc"
)

//  ===================   server request functions ===================
type rpcClient struct {
	client *rpc.Client
}

// master
func (t *rpcClient) notify(msg, key, value string) error {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("Master.Notify", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	fmt.Println(reply.Result)

	return nil
}

//  ===================   server handler functions ===================
type KVServer int

// Terminate a server
func (*KVServer) Terminate(args *shared.Args, reply *shared.Response) error {
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
