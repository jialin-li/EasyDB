package main

import (
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
func (t *rpcClient) clientConnect(clientId int) error {
	args := &shared.Args{Value: strconv.Itoa(clientId)}
	var reply shared.Response
	err := t.client.Call("KVClient.Connect", args, &reply)
	if err != nil {
		log.Println("client error:", err)
	}
	return nil
}

func (t *rpcClient) clientDisconnect(serverId int) error {
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
		switch err.Error() {
		case shared.ERR_DEP:
			fallthrough
		case shared.ERR_KEY:
			return fmt.Sprintf("%s:%s", key, err)
		default:
			log.Println("client error:", err)
		}
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
func (t *rpcClient) serverConnect(serverId int) error {
	args := &shared.Args{Value: strconv.Itoa(serverId)}
	var reply shared.Response
	err := t.client.Call("KVServer.Connect", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
}

func (t *rpcClient) serverDisconnect(serverId int) error {
	args := &shared.Args{Value: strconv.Itoa(serverId)}
	var reply shared.Response
	err := t.client.Call("KVServer.Disconnect", args, &reply)
	if err != nil {
		log.Println("server error:", err)
	}
	return nil
}

func (t *rpcClient) printStore() string {
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
	// fmt.Printf("Master: Notify, id %v \n", args.ID)
	conn, _ := shared.Dial(shared.BasePort + args.ID)
	conns[args.ID] = &rpcClient{client: rpc.NewClient(conn)}

	// if a server is notifying us, then we connect it with every other server
	// that is not itself
	if !isClientId(args.ID) {
		for _, id := range IdMap {
			if !isClientId(id) && id != args.ID {
				conns[id].serverConnect(args.ID)
				conns[args.ID].serverConnect(id)
			}
		}
	}
	shared.Outputln("notify from:", args.ID)

	if term {
		IdMap[args.ID] = args.ID
	}

	if !term {
		wg.Done()
	}
	return nil
}
