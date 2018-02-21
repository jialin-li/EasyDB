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
func (t *rpcClient) connectClient(serverId int) error {
	args := &shared.Args{Value: strconv.Itoa(serverId)}
	var reply shared.Response
	err := t.client.Call("KVClient.Connect", args, &reply)
	if err != nil {
		log.Println("client error:", err)
	}
	return nil
}

func (t *rpcClient) disconnectClient(serverId int) error {
	args := &shared.Args{Value: strconv.Itoa(serverId)}
	var reply shared.Response
	err := t.client.Call("KVClient.Disconnect", args, &reply)
	if err != nil {
		log.Println("client error:", err)
	}
	return nil
}

func (t *rpcClient) put(key, value string) error {
	args := &shared.Args{Key: key, Value: value}
	var reply shared.Response
	err := t.client.Call("KVClient.Put", args, &reply)
	if err != nil {
		log.Println("client error:", err)
	}
	return nil
}

func (t *rpcClient) get(key string) string {
	args := &shared.Args{Key: key}
	var reply shared.Response
	err := t.client.Call("KVClient.Get", args, &reply)
	if err != nil {
		log.Println("client error:", err)
	}
	return reply.Result
}

// server
func (t *rpcClient) connectServer(serverId int) error {
	args := &shared.Args{Value: strconv.Itoa(serverId)}
	var reply shared.Response
	err := t.client.Call("KVServer.Connect", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
}

func (t *rpcClient) disconnectServer(serverId int) error {
	args := &shared.Args{Value: strconv.Itoa(serverId)}
	var reply shared.Response
	err := t.client.Call("KVServer.Disconnect", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
}

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
		// if we are running in term mode then internal ids are not updated by
		// joinClient command
		if term {
			clientIds[args.ID] = args.ID
		}

	case shared.ServerType:
		conn, _ := shared.Dial(shared.ServerPort + args.ID)
		serverCalls[args.ID] = &rpcClient{client: rpc.NewClient(conn)}
		// if we are running in term mode then internal ids are not updated by
		// joinServer command
		if term {
			serverIds[args.ID] = args.ID
		}

		for _, serverId := range serverIds {
			if serverId != args.ID {
				serverCalls[args.ID].connectServer(serverId)
				serverCalls[serverId].connectServer(args.ID)
			}
		}

	default:
		log.Println("Notify failed")
		return errors.New("Unknown rpc client type")
	}

	if !term {
		wg.Done()
	}
	return nil
}
