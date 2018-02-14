package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strings"

	"github.com/jialin-li/EasyDB/shared"
)

type KVClient struct {
	client *rpc.Client
}

func (t *KVClient) Terminate(args *shared.Args, reply *shared.Response) error {
	return nil
}

func (t *KVClient) Connect(msg string) shared.Response {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, "key", "value"}
	var reply shared.Response
	err := t.client.Call("Server.Connect", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return reply
}

// Disconnect if we are clients of any other servers
func (t *KVClient) Disconnect(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Put a KV pair
func (t *KVClient) Put(args *shared.Args, reply *shared.Response) error {
	// called by the master, will issue request to server

	err := t.client.Call("KVServer.Put", args, reply)
	if err != nil {
		fmt.Printf("kvserver error: %e \n", err)
	}
	return nil
}

// Get a Value based on a key
func (t *KVClient) Get(args *shared.Args, reply *shared.Response) error {
	return nil
}
