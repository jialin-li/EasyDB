package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"strconv"

	"github.com/jialin-li/EasyDB/shared"
)

//  ===================   master request functions ===================
type rpcClient struct {
	client *rpc.Client
}

// client
func (t *rpcClient) connect(msg, key string, serverId int) error {
	fmt.Println("sending", msg)
	args := &shared.Args{Msg: msg, Key: key, Value: strconv.Itoa(serverId)}
	//args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("KVClient.Connect", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
}

func (t *rpcClient) put(key, value string) error {
	args := &shared.Args{Key: key, Value: value}
	var reply shared.Response
	err := t.client.Call("KVClient.Put", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
}

func (t *rpcClient) get(key string) string {
	args := &shared.Args{Key: key}
	var reply shared.Response
	err := t.client.Call("KVClient.Get", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return reply.Result
}

func (t *rpcClient) kill() error {
	// args := &shared.Args{Key: key, Value: value}
	// var reply shared.Response
	err := t.client.Call("KVServer.Terminate", &shared.Args{}, &shared.Response{})
	if err != nil && err != io.ErrUnexpectedEOF {
		log.Println("server error:", err)
	}
	return nil
}

// server
func (t *rpcClient) printStore() string {
	fmt.Println("sending printStore")
	args := &shared.Args{}
	var reply shared.Response
	err := t.client.Call("KVServer.DumpStore", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return reply.Result
}

//  ===================   master handler functions ===================
type Master int

func (*Master) Notify(args *shared.NotifyArgs, reply *shared.Response) error {
	// Dial back
	switch args.Type {
	case shared.ClientType:
		conn, _ := shared.Dial(shared.ClientPort + args.ID)
		clientCalls[args.ID] = &rpcClient{client: rpc.NewClient(conn)}

	case shared.ServerType:
		conn, _ := shared.Dial(shared.ServerPort + args.ID)
		serverCalls[args.ID] = &rpcClient{client: rpc.NewClient(conn)}

	default:
		log.Println("Notify failed")
		return errors.New("Unknown rpc client type")
	}

	if !term {
		wg.Done()
	}
	return nil
}
