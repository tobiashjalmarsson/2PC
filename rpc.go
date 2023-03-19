package main

import (
	"errors"
	"net"
	"net/rpc"
)

type RPCArgs struct {
	StringArg string
	BoolArg   bool
	IntArg    int
}

type RPCRes struct {
	StringResult string
	ArrayResult  []string
	BoolResult   bool
}

// ChordRPC invokes method on the node of given address.
func (node *Node) CallRPC(address string, methodName string, args *RPCArgs, res *RPCRes) error {
	if address != node.address {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return err
		}
		client := rpc.NewClient(conn)
		defer client.Close()
		err = client.Call("Node."+methodName+"RPC", args, &res)
		return err
	} else {
		return errors.New(node.address + " calls its method " + methodName + " via RPC")
	}
}

func (node *Node) GetPeersRPC(_ *RPCArgs, res *RPCRes) error {
	res.ArrayResult = node.peers
	return nil
}
