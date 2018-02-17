package main

import (
	"errors"
	"fmt"
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
	args := &shared.Args{msg, key, strconv.Itoa(serverId)}
	//args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("KVClient.Connect", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
}

func (t *rpcClient) put(msg, key, value string) error {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("KVClient.Put", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
}

func (t *rpcClient) get(msg, key, value string) error {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("KVClient.Get", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
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
