package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"
)

type Credentials struct {
	username string
	password string
}

func main() {
	fmt.Println("Getting flags")
	port := flag.String("port", "8080", "portnumber for initialized node")
	address := flag.String("ip", "localhost", "address for the initialized node")
	flag.Parse()
	fmt.Printf("Address for node is: %s:%s \n", *address, *port)

	// Get credentials
	//creds := getCredentials()
	node := CreateNode(*port, *address)
	//node.StartListener()
	// Set up RPC server
	server := rpc.NewServer()
	err := server.Register(node)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		listener, err_ := net.Listen("tcp", node.address+":"+node.port)
		if err_ != nil {
			log.Fatalln(err_)
		}
		for {
			conn, err__ := listener.Accept()
			if err__ != nil {
				log.Fatalln(err__)
			}
			// Serve RPC calls
			go server.ServeConn(conn)
		}
	}()

	timerfft := time.NewTimer(time.Duration(10) * 1000000)
	go func() {
		for {
			<-timerfft.C
			node.UpdatePeers()
			timerfft.Reset(time.Duration(10) * 1000000)
		}
	}()
	node.GetCommands()

	//fmt.Printf("Username: %s, Password: %s \n", creds.username, creds.password)
}

func getCredentials() Credentials {
	var creds Credentials
	fmt.Print("Username:")
	_, err := fmt.Scanln(&creds.username)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Password:")
	_, err = fmt.Scanln(&creds.password)
	if err != nil {
		log.Fatal(err)
	}
	return creds
}
