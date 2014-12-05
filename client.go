package poddl

import (
	"fmt"
	"github.com/agl/xmpp"
	"strings"
)

// Client is
type Client struct {
	config         *Config
	connection     *xmpp.Conn
	conversation   *Conversation
	quitchan       chan (bool)
	stanzas        chan xmpp.Stanza
	MessageHandler func(c *Client, m *xmpp.ClientMessage)
}

func (c *Client) handleIn() {
	for {
		s, err := c.connection.Next()

		if err != nil {
			fmt.Errorf("Received an error: %s", err.Error())
		}
		c.stanzas <- s
	}
}

func (c *Client) handleOut() {
	for {
		msg, ok := <-c.conversation.out
		if !ok {
			break
		}
		c.connection.Send(c.config.Contact, msg)
	}
}

// Mainloop is the loop of a Client that handles messages.
func (c *Client) Mainloop() {
	c.connection.SendPresence("", "", "")
	c.SetupFeeds()
	go c.handleIn()
	go c.handleOut()
Mainloop:
	for {
		select {
		case s := <-c.stanzas:
			switch stanza := s.Value.(type) {
			case *xmpp.ClientMessage:
				c.MessageHandler(c, stanza)
			default:
				fmt.Printf("stanza name: %s, value: %s\n", s.Name, s.Value)
			}
		case <-c.quitchan:
			close(c.conversation.in)
			close(c.conversation.out)
			break Mainloop
		}
	}
}

// NewClient creates a new Client with a connection to an XMPP server and a
// default message handler
func NewClient() (c *Client, err error) {
	poddlconf, err := readConfig()

	if err != nil {
		return nil, err
	}

	conn, err := xmpp.Dial(
		poddlconf.Address,
		poddlconf.User,
		poddlconf.Domain,
		poddlconf.Password,
		&xmpp.Config{})

	if err != nil {
		return
	}

	c = &Client{
		MessageHandler: handleClientMessage,
		config:         poddlconf,
		connection:     conn,
		conversation:   NewConversation(),
		quitchan:       make(chan bool, 1),
		stanzas:        make(chan xmpp.Stanza, 1)}
	return
}

func handleClientMessage(c *Client, m *xmpp.ClientMessage) {
	if "quit" == strings.TrimSpace(m.Body) {
		c.quitchan <- true
		return
	}
	c.conversation.in <- m.Body
}
