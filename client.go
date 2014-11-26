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
	quitchan       chan (bool)
	MessageHandler func(c *Client, m *xmpp.ClientMessage)
}

// Mainloop is the loop of a Client that handles messages.
func (c *Client) Mainloop() {
	c.connection.SendPresence("", "", "")
Mainloop:
	for {
		select {
		case <-c.quitchan:
			break Mainloop
		default:
			s, err := c.connection.Next()
			if err != nil {
				fmt.Errorf("Received an error: %s", err.Error())
			}
			switch stanza := s.Value.(type) {
			case *xmpp.ClientMessage:
				c.MessageHandler(c, stanza)
			default:
				fmt.Printf("stanza name: %s, value: %s\n", s.Name, s.Value)
			}
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
		quitchan:       make(chan bool, 1)}
	return
}

func handleClientMessage(c *Client, m *xmpp.ClientMessage) {
	fmt.Printf("Received a message from %s\n", m.From)
	fmt.Printf("It's: %s.\n", m.Body)
	if "quit" == strings.TrimSpace(m.Body) {
		c.quitchan <- true
	}
}
