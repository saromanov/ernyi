package main

import (
	"fmt"
	"github.com/saromanov/ernyi/agent"
	"github.com/saromanov/ernyi/utils"
	"github.com/saromanov/ernyi/structs"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/rpc"
)

var (
	defaultAddr = "127.0.0.1"
	rpcdefault = "127.0.0.1:9652"
	versionNum = "0.1"
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
	version = "version"
)


func CreateErnyi() {
	addrRPC := defaultAddr + ":" + utils.GenRandomPort()
	agent.CreateAgent(*name, *addr, addrRPC)
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

	Output("Join", fmt.Sprintf("Joining client by address %s", *addr))
}

func Members() {

	client, err := rpc.DialHTTP("tcp", *rpcaddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply bool
	members := &structs.MembersResponse{}
	err = client.Call("Agent.Members", &members, &reply)
	if err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
		return
	}

	if !reply {
		log.Fatal("Replay from command Join is false")
		return
	}

	fmt.Println(reply)
	fmt.Println(members)
}

func Version() {
	fmt.Println(versionNum)
}

func ProcessCommands() {
	switch *command {
	case create:
		CreateErnyi()
	case join:
		Join()
	case members:
		Members()
	case version:
		Version()
	default:
		fmt.Println("Unknown command")
	}
}

func Output(command, msg string){
	fmt.Println(fmt.Sprintf("Command: %s\nOutput: %s", command, msg))
}

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	ProcessCommands()
}
