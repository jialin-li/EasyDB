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

	"github.com/jialin-li/EasyDB/shared"
)

type rpcClient struct {
	client *rpc.Client
}

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

type KVClient int

func (*KVClient) Connect(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Disconnect if we are clients of any other servers
func (*KVClient) Disconnect(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Put a KV pair
func (*KVClient) Put(args *shared.Args, reply *shared.Response) error {
	// called by the master, will issue request to server

	//err := r.client.Call("KVServer.Put", args, reply)
	//if err != nil {
	//fmt.Printf("kvserver error: %e \n", err)
	//}
	return nil
}

// Get a Value based on a key
func (t *KVClient) Get(args *shared.Args, reply *shared.Response) error {
	return nil
}
