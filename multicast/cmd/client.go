package main

import (
	"log"

	mc "github.com/ikascrew/core/multicast"
)

func main() {

	c, err := mc.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	m, err := c.Find()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(m)
}
