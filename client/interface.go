package main

import (
	"fmt"
	"github.com/jialin-li/EasyDB/shared"
	"log"
	"net/rpc"
)

type Client struct {
	client *rpc.Client
}

func (t *Client) Terminate() error {
	return nil
}

func (t *Client) Connect(msg string) shared.Response {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, "key", "value"}
	var reply shared.Response
	err := t.client.Call("Server.Connect", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return reply
}

// Terminate a server
// Ask the server to connect to another server

// Send back the current key value store
func (t *Client) DumpStore(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Disconnect if we are clients of any other servers
func (t *Client) Disconnect(args *shared.Args, reply *shared.Response) error {
	return nil
}

// No more new write requests will be send until the next write
func (t *Client) Stabilize(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Put a KV pair
func (t *Client) Put(args *shared.Args, reply *shared.Response) error {
	return nil
}

// Get a Value based on a key
func (t *Client) Get(args *shared.Args, reply *shared.Response) error {
	return nil
}
