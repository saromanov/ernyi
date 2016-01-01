package agent

import (
	"log"
	"net"
	"fmt"
	"errors"
	"strconv"
	"net/http"
	"net/rpc"

	"github.com/saromanov/ernyi/ernyi"
	"github.com/saromanov/ernyi/structs"
	"github.com/hashicorp/memberlist"
)

// RPC agent

// Agent provides entry point for ernyi
type Agent struct {
	Ern *ernyi.Ernyi
	Tags map[string][]string
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
	agent.Tags = map[string][]string{}
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

func (agent *Agent) Members(members *[]*memberlist.Node, reply *bool) error {
	result := agent.Ern.Members()
	*members = result
	*reply = true
	return nil
}

func (agent *Agent) SetTag(item *structs.RPCSetTag, reply *bool) error {
	// Must be global update for tags
	if item == nil {
		return errors.New("Empty struct RPCSetTag")
	}

	var found bool
	for _, member := range agent.Ern.Members() {
		if item.Name == member.Name {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf(fmt.Sprintf("Node with the name %s", item.Name))
	}

	_, ok := agent.Tags[item.Tag]
	if !ok {
		agent.Tags[item.Tag] = []string{item.Name}
	} else {
		agent.Tags[item.Tag] = append(agent.Tags[item.Tag], item.Name)
	}
	return nil
}

func setupRPC(addr string) {
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
