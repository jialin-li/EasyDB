package main

import (
	"log"
	"net/rpc"
	"strconv"

	"github.com/jialin-li/EasyDB/shared"
)

//  ===================   client request functions ===================
type rpcClient struct {
	client *rpc.Client
}

// master rpc
func (t *rpcClient) notify(id int) error {
	// fmt.Println("About to notify master that client %s is up", id)
	args := &shared.NotifyArgs{Type: shared.ClientType, ID: id}
	err := t.client.Call("Master.Notify", args, nil)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return err
}

// server rpc
func (t *rpcClient) put(key, value string) error {
	args := &shared.Args{Key: key, Value: value}
	if v, ok := keyTimes[key]; ok {
		args.Time = *v
	}
	reply := &shared.Response{}
	err := t.client.Call("KVServer.Put", args, reply)
	// if err != nil {
	// 	log.Fatal("server error:", err)
	// }
	if err == nil {
		// update our time and add it to the map
		incTime(&reply.Time)
		keyTimes[key] = &reply.Time
		shared.Outputf("client put: %s:%s, time %v \n", key, value, reply.Time)
	}
	// TODO: read the time stamp from reply and do things with it
	return err
}

func (t *rpcClient) get(key string, reply *shared.Response) error {
	args := &shared.Args{Key: key}
	err := t.client.Call("KVServer.Get", args, reply)
	// ?? do we need to deal with time stamp?
	// if err != nil {
	// 	log.Fatal("server error:", err)
	// }
	return err
}

//  ===================   client handler functions ===================
type KVClient int

func (*KVClient) Connect(args *shared.Args, reply *shared.Response) error {
	// server id that we are going to attempt to connect to
	serverId, err := strconv.Atoi(args.Value)
	if err != nil {
		log.Println(err)
	}
	err = setupConn(serverId)
	//fmt.Println("connected")
	//TODO: give appropriate reply if connection already exists
	return err
}

// Disconnect if we are clients of any other servers
func (*KVClient) Disconnect(args *shared.Args, reply *shared.Response) error {
	// server id that we are going to disconnect from
	serverId, err := strconv.Atoi(args.Value)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println("disconnecting from:", serverId)
	err = serverCalls[serverId].client.Close()
	// delete the connection from the map of server calls
	delete(serverCalls, serverId)
	//TODO: give appropriate reply if connection does not exist

	return err
}

// Put a KV pair
func (*KVClient) Put(args *shared.Args, reply *shared.Response) (err error) {
	// called by the master, will issue request to server
	// go through list of server connections
	for sid := range serverCalls {
		if err = serverCalls[sid].put(args.Key, args.Value); err == rpc.ErrShutdown {
			// lazily remove closed server connection
			delete(serverCalls, sid)
		} else if err == nil {
			return nil
		}
	}
	return err
}

// Get a Value based on a key
func (t *KVClient) Get(args *shared.Args, reply *shared.Response) (err error) {
	for sid := range serverCalls {
		if err = serverCalls[sid].get(args.Key, reply); err == rpc.ErrShutdown {
			// lazily remove closed server connection
			delete(serverCalls, sid)
		} else if err == nil {
			return nil
		}
	}
	return err
}

func setupConn(serverId int) error {
	// connect to the specified server
	conn, err := shared.Dial(serverId + shared.BasePort)
	// overwrite existing connection, if any
	serverCalls[serverId] = &rpcClient{client: rpc.NewClient(conn)}

	return err
}

func incTime(time *shared.Time) {
	if clientId >= shared.MaxClient || clientId < 0 {
		log.Println("clientId exceeds the slot given in the current timestamp")
		return
	}
	time.Clock[clientId-shared.ClientStart]++
}
