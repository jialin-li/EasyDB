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

func (t *rpcClient) get(msg, key, value string) error {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("KVServer.Get", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return nil
}

//  ===================   client handler functions ===================
type KVClient int

func (*KVClient) Connect(args *shared.Args, reply *shared.Response) error {
	// server id that we are going to attempt to connect to
	serverId, err := strconv.Atoi(args.Value)
	if err != nil {
		log.Println(err)
	}

	conn, _ := shared.Dial(shared.ServerPort + serverId)

	serverCalls[serverId] = &rpcClient{client: rpc.NewClient(conn)}
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

	serverCalls[1].put(args.Msg, args.Key, args.Value)

	//err := r.client.Call("KVServer.Put", args, reply)
	//if err != nil {
	//fmt.Printf("kvserver error: %e \n", err)
	//}
	return nil
}

// Get a Value based on a key
func (t *KVClient) Get(args *shared.Args, reply *shared.Response) error {
	serverCalls[1].get(args.Msg, args.Key, args.Value)
	return nil
}
