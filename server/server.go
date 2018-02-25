package main

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
	"sync"

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
	return bulkLoad()
}

// Put a KV pair
func (*KVServer) Put(args *shared.Args, reply *shared.Response) error {
	return put(args, reply)
}

// Get a Value based on a key
func (*KVServer) Get(args *shared.Args, reply *shared.Response) error {
	return get(args, reply)
}

// Load a series of key:value pairs along with their timestamp
func (*KVServer) BulkLoad(args *shared.Args, reply *shared.Response) error {
	return bulkSave(args, reply)
}

// actual kv functions
func put(args *shared.Args, reply *shared.Response) error {
	if v, ok := db[args.Key]; ok {
		v.value = args.Value
		//shared.Outputf("prev server: time %v \n", v.time)
		v.time.Update(&args.Time)
	} else {
		db[args.Key] = &dbValue{value: args.Value, time: args.Time}
	}
	// increment server's time by 1
	incTime(&db[args.Key].time)
	shared.Outputf("server: time:%v, id:%v\n", db[args.Key].time, serverId)
	reply.Time = db[args.Key].time

	return nil
}

func get(args *shared.Args, reply *shared.Response) error {
	var val string
	// check if key is present in server's db
	if v, ok := db[args.Key]; ok {
		// if key exists but the client has a later timestamp for that key,
		// then return DepError
		if args.Time.IsLaterThan(&db[args.Key].time) {
			return &shared.DepError{}
		}
		// otherwise we return our value with it's newer timestamp
		val = v.value
		reply.Time = v.time
	} else {
		return &shared.KeyError{}
	}
	shared.Outputf("server: time:%v, id:%v\n", db[args.Key].time, serverId)

	reply.Result = fmt.Sprintf("%s:%s", args.Key, val)
	return nil
}

func bulkLoad() error {
	var stabilizeWait sync.WaitGroup
	// TODO: gather values and serialize
	args := &shared.Args{Key: strconv.Itoa(serverId)}

	// send our entire database to every connected server
	for sid := range serverCalls {
		log.Printf("sending bulk to server: %v from %v", sid, serverId)
		stabilizeWait.Add(1)
		// make rpc calls in new go routines
		go func(id int) {
			err := serverCalls[id].client.Call("KVServer.BulkLoad", args, nil)
			if err == rpc.ErrShutdown {
				delete(serverCalls, id)
			}
			stabilizeWait.Done()
		}(sid)
	}
	// wait for all BulkLoad rpc calls to finish
	stabilizeWait.Wait()
	return nil
}

func bulkSave(args *shared.Args, reply *shared.Response) error {
	shared.Outputf("server %v saving data from %v\n", serverId, args.Key)
	return nil
}

func incTime(time *shared.Time) {
	if serverId < shared.ServerStart || serverId >= (10+shared.ServerStart) {
		log.Printf("server exceeds the slot given in the current timestamp: id %v \n", serverId)
		return
	}
	time.Clock[serverId-shared.ServerStart+shared.MaxClient]++
}
