package ernyi

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/saromanov/ernyi/ernyi/event"
	"github.com/hashicorp/memberlist"
)

// Definition of errors
var (
	errEmptyName        = errors.New("Member must contain address")
	errEmptyListMembers = errors.New("List of members is empty")
)

// Erny provides main struct
type Ernyi struct {
	// main item for member list
	mlist      *memberlist.Memberlist
	memberlock *sync.RWMutex
	tags       map[string][]string
	event      chan event.Event
	addr       string
	events     map[string][]string
	fails      map[string][]string
	success    map[string][]string
}

// CreateErnyi provides new object of ernyi
func CreateErnyi(config *Config) *Ernyi {
	defaultConfig := memberlist.DefaultLANConfig()
	bindAddr := defaultConfig.BindAddr
	bindPort := defaultConfig.BindPort
	if config != nil && config.MemberlistConfig != nil {
		bindAddr = config.MemberlistConfig.BindAddr
		bindPort = config.MemberlistConfig.BindPort
		defaultConfig = config.MemberlistConfig
	}
	ern := new(Ernyi)
	ern.memberlock = &sync.RWMutex{}
	Addr := &net.TCPAddr{
		IP:   net.ParseIP(bindAddr),
		Port: bindPort,
	}
	ern.addr = Addr.String()
	// Create basic memberlist model
	mlist, err := memberlist.Create(defaultConfig)
	if err != nil {
		log.Fatal(err)
	}
	ern.tags = map[string][]string{}
	ern.mlist = mlist
	ern.event = make(chan event.Event, 64)
	err = ern.Join(ern.addr)
	if err != nil {
		log.Fatal(err)
	}
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
		return err
	}

	if nummembers != 1 {
		return fmt.Errorf("Expected number of joining nodes %d. Found - %d", 1, nummembers)
	}

	return nil
}

// JoinMany provides joining several nodes at once
func (ern *Ernyi) JoinMany(addrs []string) error {
	lenAddr := len(addrs)
	if lenAddr == 0 {
		return errEmptyListMembers
	}
	ern.memberlock.Lock()
	defer ern.memberlock.Unlock()
	nummembers, err := ern.mlist.Join(addrs)
	if err != nil {
		return err
	}

	if lenAddr != nummembers {
		return fmt.Errorf("Expected number of joining nodes %d. Found - %d", lenAddr,
			nummembers)
	}
	return nil
}

// Leave provides closes for ernyi
func (ern *Ernyi) Leave() error {
	ern.memberlock.Lock()
	defer ern.memberlock.Unlock()
	return ern.mlist.Leave(3*time.Second)
}

// Tags add tags for node
func (ern *Ernyi) Tags(nodename string, tags []string) {
	_, ok := ern.tags[nodename]
	if !ok {
		return
	}
	ern.tags[nodename] = tags
}

// Send provides sending data to nodes on TCP
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
		return time.Second, fmt.Errorf("Addrname parameter can't be empty")
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

// LocalNode returns current local node
func (ern *Ernyi) LocalNode()*memberlist.Node {
	return ern.mlist.LocalNode()
}

