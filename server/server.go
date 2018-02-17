package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/jialin-li/EasyDB/shared"
)

//  ===================   server request functions ===================
type rpcClient struct {
	client *rpc.Client
}

func (t *rpcClient) notify(id int) error {
	// fmt.Println("About to notify master that client %s is up", id)
	args := &shared.NotifyArgs{Type: shared.ClientType, ID: id}
	// var reply shared.Response
	err := t.client.Call("Master.Notify", args, nil)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return err
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
	return put(args, reply)
}

// Get a Value based on a key
func (*KVServer) Get(args *shared.Args, reply *shared.Response) error {
	return get(args, reply)
}

// actual kv functions
func put(args *shared.Args, reply *shared.Response) error {
	if v, ok := db[args.Key]; ok {
		v.value = args.Value
		// TODO: update to real time
		v.time = args.Time
	} else {
		db[args.Key] = &dbValue{value: args.Value, time: args.Time}
	}
	// fmt.Println("new db entry:")
	// fmt.Printf(db[args.Key].value)
	reply.Time = db[args.Key].time
	fmt.Printf("%s:%s \n", args.Key, args.Value)
	return nil
}

func get(args *shared.Args, reply *shared.Response) error {
	// TODO: verify time stamp
	if v, ok := db[args.Key]; ok {
		reply.Result = v.value
		reply.Time = v.time
	} else {
		reply.Result = shared.ERR_KEY
	}
	// fmt.Println("new db entry:")
	// fmt.Printf(db[args.Key].value)

	// fmt.Printf("%s:%s \n", args.Key, args.Value)
	// reply.Result = "done"
	return nil
}
