package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ikascrew/core/window"
	"github.com/ikascrew/plugin/file"
)

func main() {

	flag.Parse()
	args := flag.Args()

	if len(args) <= 0 {
		fmt.Println("arguments MP4 file")
		os.Exit(1)
	}

	v, err := file.New(args[0])
	if err != nil {
		panic(err)
	}

	win, err := window.New("Sample Window")

	err = win.Play(v)
	if err != nil {
		panic(err)
	}

}
