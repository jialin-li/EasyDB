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
func (t *rpcClient) notify(msg, key, value string) error {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("Master.Notify", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return nil
}

// server rpc
func (t *rpcClient) put(msg, key, value string) error {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("KVServer.Put", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return nil
}

func (t *rpcClient) get(key string) *shared.Response {
	fmt.Println("getting ", key)
	// TODO: clean this up
	args := &shared.Args{"", key, ""}
	reply := &shared.Response{}
	err := t.client.Call("KVServer.Get", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return reply
}

//  ===================   client handler functions ===================
type KVClient int

func (*KVClient) Connect(args *shared.Args, reply *shared.Response) error {
	// server id that we are going to attempt to connect to
	serverId, err := strconv.Atoi(args.Value)
	if err != nil {
		log.Println(err)
	}
	connectServer(serverId)
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
		serverCalls[sid].put(args.Msg, args.Key, args.Value)
		return nil
	}

	//err := r.client.Call("KVServer.Put", args, reply)
	//if err != nil {
	//fmt.Printf("kvserver error: %e \n", err)
	//}
	return nil
}

// Get a Value based on a key
func (t *KVClient) Get(args *shared.Args, reply *shared.Response) error {
	for sid := range serverCalls {
		reply = serverCalls[sid].get(args.Key)
	}
	return nil
}

func connectServer(serverId int) {
	// connect to the specified server
	conn, _ := shared.Dial(serverId + shared.ServerPort)
	serverCalls[serverId] = &rpcClient{client: rpc.NewClient(conn)}
}
