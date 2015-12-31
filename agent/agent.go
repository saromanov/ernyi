package agent

import (
	"log"
	"net"
	"strconv"
	"net/http"
	"net/rpc"

	"github.com/saromanov/ernyi/ernyi"
	"github.com/hashicorp/memberlist"
)

// RPC agent

// Agent provides entry point for ernyi
type Agent struct {
	Ern *ernyi.Ernyi
}

func CreateAgent(name, addr, rpcaddr string) {
	agent := new(Agent)
	mconfig := memberlist.DefaultLANConfig()
	if name == "" {
		log.Fatal("Name must be non-empty")
	}

	if addr == "" {
		log.Fatal("Address must be non-empty")
	}
	shost, sport, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(err)
	}

	mconfig.Name = name
	mconfig.BindAddr = shost
	res, erratoi := strconv.Atoi(sport)
	if erratoi != nil {
		log.Fatal(erratoi)
	}
	mconfig.BindPort = res

	cfg := &ernyi.Config{
		MemberlistConfig: mconfig,
	}

	value := ernyi.CreateErnyi(cfg)
	agent.Ern = value
	rpc.Register(agent)
	setupRPC(rpcaddr)
	agent.Ern.Start()
}

func (agent *Agent) Join(addr string, reply *bool) error {
	err := agent.Ern.Join(addr)
	if err == nil {
		*reply = true
	}

	return err
}

func setupRPC(addr string) {
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
