package ernyi

import (
    "github.com/hashicorp/memberlist"
    "sync"
    "errors"
    "log"
    "fmt"
)

var (
	errEmptyName = errors.New("Member must contain address")
	errEmptyListMembers = errors.New("List of members is empty")
)

type Ernyi struct {
	mlist  *memberlist.Memberlist
	memberlock  *sync.RWMutex
	tags   map[string][]string
}

func CreateErnyi(config *Config)*Ernyi {
	ern := new(Ernyi)
	ern.memberlock = &sync.RWMutex{}
	mlist, err  = Create(config.MemberlistConfig)
	if err != nil {
		log.Fatal(err)
	}
	ern.tags = map[string][]string{}
	ern.mlist = mlist
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
	nummembers, err := ern.mlist.Join(addr)
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
	return ern.Members()
}