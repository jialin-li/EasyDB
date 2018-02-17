package main

import (
	//	"bufio"
	//	"flag"
	"fmt"
	"log"
	//	"net"
	"net/rpc"
	//	"os"
	//	"strings"
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
	args := &shared.NotifyArgs{Type: shared.ServerType, ID: id}
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
		args.Time = v
	}
	reply := &shared.Response{}
	err := t.client.Call("KVServer.Put", args, reply)
	// if err != nil {
	// 	log.Fatal("server error:", err)
	// }
	// TODO: read the time stamp from reply and do things with it
	return err
}

func (t *rpcClient) get(key string, reply *shared.Response) error {
	fmt.Printf("getting %v \n", key)
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
	setupConn(serverId)
	fmt.Println("connected")
	//TODO: give appropriate reply if connection already exists
	return err
}

// Disconnect if we are clients of any other servers
func (*KVClient) Disconnect(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Put a KV pair
func (*KVClient) Put(args *shared.Args, reply *shared.Response) error {
	// called by the master, will issue request to server
	// go through list of server connections
	for sid := range serverCalls {
		return serverCalls[sid].put(args.Key, args.Value)
	}
	return nil
}

// Get a Value based on a key
func (t *KVClient) Get(args *shared.Args, reply *shared.Response) error {
	for sid := range serverCalls {
		return serverCalls[sid].get(args.Key, reply)
	}
	return nil
}

func setupConn(serverId int) {
	// connect to the specified server
	conn, _ := shared.Dial(serverId + shared.ServerPort)
	serverCalls[serverId] = &rpcClient{client: rpc.NewClient(conn)}
}
