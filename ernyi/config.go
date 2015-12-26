package ernyi

import (
	"github.com/hashicorp/memberlist"
)

type Config struct {
	MemberlistConfig *memberlist.Config
	// Addr is address for server
	Addr string
}
