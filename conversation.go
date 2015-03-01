package poddl

import (
	"sync"
)

type Conversation struct {
	lock sync.Mutex
	in   chan string
	out  chan string
}

func NewConversation () (c *Conversation) {
	c = &Conversation{in: make(chan string), out: make(chan string)}
	return
}
