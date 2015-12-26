package ernyi

import (
    "github.com/hashicorp/memberlist"
    "sync"
    "errors"
    "log"
    "fmt"
    "./event"
)

var (
	errEmptyName = errors.New("Member must contain address")
	errEmptyListMembers = errors.New("List of members is empty")
)

type Ernyi struct {
	mlist  *memberlist.Memberlist
	memberlock  *sync.RWMutex
	tags   map[string][]string
	event  chan event.Event
	addr string
}

func CreateErnyi(config *Config)*Ernyi {
	ern := new(Ernyi)
	ern.memberlock = &sync.RWMutex{}
	mlist, err := memberlist.Create(config.MemberlistConfig)
	if err != nil {
		log.Fatal(err)
	}
	ern.tags = map[string][]string{}
	ern.mlist = mlist
	ern.event = make(chan event.Event,64)
	ern.addr = config.Addr
	return ern
}

// Join provides joining of the new member
func (ern *Ernyi) Join(addr string) error{
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

// Tags add tags for node
func (ern *Ernyi) Tags(nodename string, tags[]string) {
	ern.tags[nodename] = tags
}

func (ern *Ernyi) Send(addr string, msg []byte) error {
	node := ern.mlist.LocalNode()

	if node == nil {
		return fmt.Errorf("Can't get local node")
	}
	fmt.Println("NUm MEMBERS: ", ern.mlist.NumMembers())
	ern.mlist.SendToTCP(node, msg)

	return nil
}

// Info returns information about Ernyi
func (ern *Ernyi) Info() map[string] string {
	return map[string] string {
		"protocol_version": fmt.Sprintf("%s", ern.mlist.ProtocolVersion()),
	}
}


// Start provides basic start if Ernyi
func (ern *Ernyi) Start() {
	go func(){
		for {
			select {
				case item := <- ern.event:
					eventname := item.String()
					if eventname == "stop" {
						return
					}
			}
		}
	}()

	StartServer(ern.addr)
}

// Stop provides stopping of Ernyi
func (ern *Ernyi) Stop() error {
	var err error
	 err = ern.mlist.Shutdown()
	 if err != nil {
	 	return err
	 }

	 return nil
}

// Members return current alive members
func (ern *Ernyi) Members()[]*memberlist.Node {
	return ern.mlist.Members()
}