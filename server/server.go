package main

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"

	"github.com/jialin-li/EasyDB/shared"
)

//  ===================   server request functions ===================
type rpcClient struct {
	client *rpc.Client
}

func (t *rpcClient) notify(id int) error {
	// fmt.Println("About to notify master that client %s is up", id)
	args := &shared.NotifyArgs{Type: shared.ServerType, ID: id}
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
	return serverListener.Close()
}

// Ask the server to connect to another server
func (*KVServer) Connect(args *shared.Args, reply *shared.Response) error {
	serverId, err := strconv.Atoi(args.Value)
	if err != nil {
		log.Println("server error:", err)
	}
	shared.Outputln("server: connected to server")

	conn, _ := shared.Dial(shared.BasePort + serverId)
	serverCalls[serverId] = &rpcClient{client: rpc.NewClient(conn)}
	return nil
}

// Send back the current key value store
func (*KVServer) DumpStore(args *shared.Args, reply *shared.Response) error {
	store := dumpStore()
	reply.Result = store
	return nil
}

// Disconnect if we are clients of any other servers
func (*KVServer) Disconnect(args *shared.Args, reply *shared.Response) error {
	// server id that we are going to disconnect from
	serverId, err := strconv.Atoi(args.Value)
	if err != nil {
		log.Println(err)
	}
	shared.Outputln("server: disconnected from server")

	err = serverCalls[serverId].client.Close()
	// delete the connection from the map of server calls
	delete(serverCalls, serverId)
	//TODO: give appropriate reply if connection does not exist

	return err
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
	// increment client's time by 1
	incTime(&args.Time)
	shared.Outputf("server: time %v \n", args.Time)
	if v, ok := db[args.Key]; ok {
		v.value = args.Value
		// TODO: update to real time
		// NOTE: right now just incrementing 1 to the time provided by the client
		v.time = args.Time
	} else {
		db[args.Key] = &dbValue{value: args.Value, time: args.Time}
	}
	// fmt.Println("new db entry:")
	// fmt.Printf(db[args.Key].value)
	reply.Time = db[args.Key].time
	// fmt.Printf("%s:%s \n", args.Key, args.Value)
	return nil
}

func get(args *shared.Args, reply *shared.Response) error {
	var val string
	// TODO: verify time stamp
	if v, ok := db[args.Key]; ok {
		val = v.value
		reply.Time = v.time
	} else {
		val = shared.ERR_KEY
	}
	// fmt.Println("new db entry:")
	// fmt.Printf(db[args.Key].value)

	// fmt.Printf("%s:%s \n", args.Key, args.Value)
	// reply.Result = "done"
	reply.Result = fmt.Sprintf("%s:%s", args.Key, val)
	return nil
}

func incTime(time *shared.Time) {
	if serverId < shared.ServerStart || serverId >= (10+shared.ServerStart) {
		log.Printf("server exceeds the slot given in the current timestamp: id %v \n", serverId)
		return
	}
	time.Clock[serverId-shared.ServerStart+shared.MaxClient]++
}
