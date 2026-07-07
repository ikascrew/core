package main

import (
	"log"

	mc "github.com/ikascrew/core/multicast"
)

func main() {

	s, err := mc.NewServer(
		mc.ServerName("ikasbox"),
		mc.Type(mc.TypeIkasbox),
	)

	if err != nil {
		log.Fatal(err)
	}

	err = s.Dial()
	if err != nil {
		log.Fatal(err)
	}

}
