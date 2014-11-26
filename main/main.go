package main

import (
	"fmt"
	"github.com/mineo/poddl"
)

//
func main() {
	client, err := poddl.NewClient()
	if err != nil {
		fmt.Println(err.Error())
	}
	client.Mainloop()
}
