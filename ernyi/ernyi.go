package ernyi

import (
    "github.com/hashicorp/memberlist"
    "sync"
    "errors"
    "log"
)

var (
	errEmptyName = errors.New("Member must contain address")
)

type Ernyi struct {
	mlist  *memberlist.Memberlist
	memberlock  *sync.RWMutex
}

func CreateErnyi(config *Config)*Ernyi {
	ern := new(Ernyi)
	ern.memberlock = &sync.RWMutex{}
	mlist, err  = Create(config.MemberlistConfig)
	if err != nil {
		log.Fatal(err)
	}
	ern.mlist = mlist
	return ern
}

// Joib provides joining of the new member
func (ern *Ernyi) Join(addr string) error{
	if addr == "" {
		return errEmptyName
	}

	ern.mlist.Join([]string{addr})
}

func (ern *Ernyi) Members()[]*memberlist.Node {
	return ern.Members()
}