package main

import (
	"fmt"
	"github.com/jialin-li/EasyDB/shared"
	"log"
	"net"
	"net/rpc"
	"strconv"
)

//  ===================   master request functions ===================
type rpcClient struct {
	client *rpc.Client
}

// client
func (t *rpcClient) put(msg, key, value string) error {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("KVClient.Put", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return nil
}

func (t *rpcClient) get(msg, key, value string) error {
	fmt.Println("sending", msg)
	args := &shared.Args{msg, key, value}
	var reply shared.Response
	err := t.client.Call("KVClient.Get", args, &reply)
	if err != nil {
		log.Fatal("server error:", err)
	}
	return nil
}

//  ===================   master handler functions ===================
type Master int

func (*Master) Notify(args *shared.Args, reply *shared.Response) error {
	*reply = shared.Response{"it worked"}
	fmt.Println(args.Msg, args.Key, args.Value)
	id, err := strconv.Atoi(args.Value)
	if err != nil {
		log.Println(err)
		return err
	}

	// Dial back
	switch t, _ := strconv.Atoi(args.Key); t {
	case shared.ClientType:
		port := shared.ClientPort + id
		conn, err := net.Dial("tcp",
			"localhost:"+strconv.Itoa(port))
		if err != nil {
			fmt.Printf("Connection: %e \n", err)
		}
		// add client connection to map
		clientConnections[id] = connection{shared.ClientType, port, conn}

		remoteCall = &rpcClient{client: rpc.NewClient(conn)}

	case shared.ServerType:
		port := shared.ServerPort + id
		conn, err := net.Dial("tcp",
			"localhost:"+strconv.Itoa(port))
		if err != nil {
			fmt.Printf("Connection: %e \n", err)
		}
		// add server connection to map
		serverConnections[id] = connection{shared.ServerType, port, conn}
	default:
		log.Println("Notify failed")
	}
	return nil
}
