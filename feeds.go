package poddl

import (
	"fmt"
	"github.com/SlyMarbo/rss"
	"time"
)

type FeedHandler struct {
	conv   *Conversation
	ticker *time.Ticker
	url    string
	In     chan string
	Out    chan string
}

func (f *FeedHandler) start() {
	feed, err := rss.Fetch(f.url)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _ = range f.ticker.C {
		f.conv.lock.Lock()
		f.conv.out <- fmt.Sprintf("Getting %s", f.url)
		err = feed.Update()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if feed.Unread == 0 {
			f.conv.lock.Unlock()
			continue
		}

		for _, item := range feed.Items {
			if item.Read {
				continue
			}

			f.conv.out <- fmt.Sprintf("Do you want me to download %s (%s)", item.Title, item.Link)
			fmt.Println("waiting for a response")
			_ = <- f.conv.in
			item.Read = true
		}
		f.conv.lock.Unlock()
		break
	}
}

func (f *FeedHandler) stop() {
	f.ticker.Stop()
}

func (c *Client) SetupFeeds() {
	for _, feed := range c.config.Feeds {
		// t := time.NewTicker(24 * time.Hour)
		t := time.NewTicker(10 * time.Second)
		handler := FeedHandler{
			conv:   c.conversation,
			url:    feed.URL,
			ticker: t}
		go handler.start()
	}
}
