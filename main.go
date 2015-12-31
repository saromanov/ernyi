package main

import (
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/saromanov/ernyi/agent"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/rpc"
)

var (
	rpcdefault = "127.0.0.1:9652"
)

var (
	command = kingpin.Arg("command", "Command").Required().String()
	name    = kingpin.Flag("name", "Name of the node").String()
	addr    = kingpin.Flag("addr", "Address of Ernyi node in format host:port").String()
	rpcaddr = kingpin.Flag("rpcaddr", "RPC address").Default(rpcdefault).String()
)

var (
	create  = "create"
	join    = "join"
	info    = "info"
	members = "members"
)


func CreateErnyi() {
	agent.CreateAgent(*name, *addr, *rpcaddr)
}

func Join() {
	if *rpcaddr == "" {
		log.Fatal("RPC address is empty")
	}

	client, err := rpc.DialHTTP("tcp", *rpcaddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply bool
	err = client.Call("Agent.Join", *addr, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}

	if !reply {
		log.Fatal("Replay from command Join is false")
	}
}

func Members() {
	client, err := rpc.DialHTTP("tcp", *rpcaddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply bool
	var members []*memberlist.Node
	err = client.Call("Agent.Members", &members, &reply)
	if err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}

	if !reply {
		log.Fatal("Replay from command Join is false")
	}

	fmt.Println(members)
}

func ProcessCommands() {
	switch *command {
	case create:
		CreateErnyi()
	case join:
		Join()
	case members:
		Members()
	default:
		fmt.Println("Unknown command")
	}
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	ProcessCommands()
}
