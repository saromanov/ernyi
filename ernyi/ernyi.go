package ernyi

import (
	"github.com/saromanov/ernyi/ernyi/event"
	"errors"
	"fmt"
	"github.com/hashicorp/memberlist"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"string"
)

var (
	errEmptyName        = errors.New("Member must contain address")
	errEmptyListMembers = errors.New("List of members is empty")
)

type Ernyi struct {
	mlist      *memberlist.Memberlist
	memberlock *sync.RWMutex
	tags       map[string][]string
	event      chan event.Event
	addr       string
}

// CreateErnyi provides new object of ernyi
func CreateErnyi(config *Config) *Ernyi {
	ern := new(Ernyi)
	ern.memberlock = &sync.RWMutex{}
	Addr := &net.TCPAddr{
		IP:   net.ParseIP(config.MemberlistConfig.BindAddr),
		Port: config.MemberlistConfig.BindPort,
	}
	ern.addr = Addr.String()
	mlist, err := memberlist.Create(config.MemberlistConfig)
	if err != nil {
		log.Fatal(err)
	}
	ern.tags = map[string][]string{}
	ern.mlist = mlist
	ern.event = make(chan event.Event, 64)
	ern.Join(ern.addr)
	return ern
}

// Join provides joining of the new member
func (ern *Ernyi) Join(addr string) error {
	if addr == "" {
		return errEmptyName
	}
	ern.memberlock.Lock()
	defer ern.memberlock.Unlock()
	nummembers, err := ern.mlist.Join([]string{addr})
	if err != nil {
		fmt.Println("ERRR: ", nummembers, err)
		return err
	}

	if nummembers != 1 {
		return fmt.Errorf("Expected number of joining nodes %d. Found - %d", 1, nummembers)
	}

	return nil
}

// JoinMany provides joining several nodes at once
func (ern *Ernyi) JoinMany(addrs []string) error {
	if len(addrs) == 0 {
		return errEmptyListMembers
	}
	ern.memberlock.Lock()
	defer ern.memberlock.Unlock()
	nummembers, err := ern.mlist.Join(addrs)
	if err != nil {
		return err
	}

	if len(addrs) != nummembers {
		return fmt.Errorf("Expected number of joining nodes %d. Found - %d", len(addrs),
			nummembers)
	}
	return nil
}

// Leave provides closes for ernyi
func (ern *Ernyi) Leave() error {
	ern.memberlock.Lock()
	defer ern.memberlock.Unlock()
	return ern.mlist.Leave()
}

// Tags add tags for node
func (ern *Ernyi) Tags(nodename string, tags []string) {
	ern.tags[nodename] = tags
}

func (ern *Ernyi) Send(addr string, msg []byte) error {
	node := ern.mlist.LocalNode()

	if node == nil {
		return fmt.Errorf("Can't get local node")
	}
	ern.mlist.SendToTCP(node, msg)

	return nil
}

// Ping provides ping to the node with the specified name
func (ern *Ernyi) Ping(addrname string, addr net.Addr) (time.Duration, error) {
	if addrname == "" {
		return time.Duration{}, fmt.Errorf("Addrname can't be empty")
	}
	return ern.mlist.Ping(addrname, addr)
}

// Info returns information about Ernyi
func (ern *Ernyi) Info() map[string]string {
	return map[string]string{
		"protocol_version": fmt.Sprintf("%s", ern.mlist.ProtocolVersion()),
		"time": time.Now().String(),
	}
}

// Start provides basic start if Ernyi
func (ern *Ernyi) Start() {
	go func() {
		for {
			select {
			case item := <-ern.event:
				eventname := item.String()
				if eventname == "stop" {
					return
				}

				// ping to random hosts
				if eventname == "ping" {

				}
			}
		}
	}()

	ern.receiveExit()
}

// Check exit from ernyi
func (ern *Ernyi) receiveExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	<-c
	ern.Stop()
	os.Exit(1)
}

// Stop provides stopping of Ernyi
func (ern *Ernyi) Stop() error {
	err := ern.mlist.Shutdown()
	if err != nil {
		return err
	}

	return nil
}

// Members return current alive members
func (ern *Ernyi) Members() []*memberlist.Node {
	return ern.mlist.Members()
}
