package ernyi

import (
   "time"
)

// Member provides basic member for Ernyi
type Member struct {
	// Address of the node
	Addr string
	// Name of the node
	Name string
	// Time whn node is added
	CrtTime time.Time
}
