package structs

import (
	"github.com/hashicorp/memberlist"
)
type RPCSetTag struct {
	Tag  string
	Name string
}

type MembersResponse struct {
	Members []*memberlist.Node
}