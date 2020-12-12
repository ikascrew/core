package main

import (
	"github.com/ikascrew/core/window"
	term "github.com/ikascrew/plugin/terminal"
)

func main() {

	v, err := term.New("Terminal Sample\n Hello ikascrew!")
	if err != nil {
		panic(err)
	}

	win, err := window.New("Sample Terminal Window")

	err = win.Play(v)
	if err != nil {
		panic(err)
	}

}
