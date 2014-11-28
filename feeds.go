package poddl

import (
	"fmt"
)

func handleFeed(conv *Conversation, uri string)  {
	fmt.Printf("Handling %s\n", uri)
	in := make(chan string)
	out := make(chan string)
	for {
		conv.aquire(in, out)
		conv.release()
	}
}

func (c *Client) SetupFeeds() {
	for _, feed := range c.config.Feeds {
		fmt.Println(feed)
	}
}
