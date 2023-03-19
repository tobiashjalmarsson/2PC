package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Node struct {
	address string
	port    string
	peers   []string
}

func CreateNode(port string, address string) Node {
	var node Node
	node.address = address
	node.port = port

	return node
}

// UpdatePeers Contact peers and check if they contain peers
// that we don't have access to, if it does
// add it to local list of peers
// If we can't connect remove it from peers as it is offline
// TODO: Call periodically to keep peers updated
func (n *Node) UpdatePeers() {
	for idx, peer := range n.peers {
		var res RPCRes
		err := n.CallRPC(peer, "GetPeers", nil, &res)

		if err != nil {
			// Could not connect to peer so remove it from the list of peers
			n.peers = append(n.peers[:idx], n.peers[idx+1:]...)
		}
		for _, remotepeer := range res.ArrayResult {
			containspeer := false
			for _, oldpeer := range n.peers {
				if oldpeer == remotepeer {
					containspeer = true
				}
			}

			if !containspeer {
				n.peers = append(n.peers, remotepeer)
			}
		}
	}
}

func (n *Node) StartListener() {
	err := rpc.Register(n)
	if err != nil {
		log.Fatalf("Failed to register node for rpc")
	}

	// Initialize listener
	listener, e := net.Listen("tcp", n.port)
	if e != nil {
		log.Fatalf("Failed to setup listener")
	}

	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatalf("Error listening", err)
	}
}

// Function to parse commands from the user
func (n *Node) GetCommands() {
	for {
		var command string
		fmt.Print("Node: ")
		_, err := fmt.Scanln(&command)
		if err != nil {
			log.Fatal(err)
		}

		switch command {
		case "help":
			fmt.Println("Option 1: upload, upload new file to peers")
			fmt.Println("Option 2: download, download new files from peers")
			fmt.Println("Option 3: get-peers, update peer-list from server")
		case "upload":
			var filepath string
			fmt.Print("filepath:")
			_, err := fmt.Scanln(&filepath)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Uploading %s \n", filepath)
		case "download":
			fmt.Println("Downloading files from peers...")
		case "get-peers":
			fmt.Println("Getting peers from the server..")
		default:
			fmt.Println("Invalid command, use <help> for options")
		}
	}
}
