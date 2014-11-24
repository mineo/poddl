package main

import (
	"github.com/agl/xmpp"
	"fmt"
	"github.com/mineo/poddl"
)

//
func main() {
	cfg := &xmpp.Config{}
	client, err := poddl.NewClient()
	if err != nil {
		fmt.Println(err.Error())
	}
	client.Mainloop()
}
