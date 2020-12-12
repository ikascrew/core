package sync

import (
	"sync"
)

type Group struct {
	errors []error
	sem    chan struct{}
	*sync.WaitGroup
}

func NewGroup(num int) *Group {
	g := Group{}
	g.WaitGroup = &sync.WaitGroup{}
	if num > 0 {
		g.sem = make(chan struct{}, num)
	} else {
		g.sem = nil
	}
	g.errors = nil
	return &g
}

func (g *Group) Add() {
	if g.sem != nil {
		g.sem <- struct{}{}
	}
	g.WaitGroup.Add(1)
	return
}

func (g *Group) AddError(e error) {
	if g.errors == nil {
		g.errors = make([]error, 0)
	}
	g.errors = append(g.errors, e)
}

func (g *Group) Done() {
	if g.sem != nil {
		<-g.sem
	}
	g.WaitGroup.Done()
	return
}

func (g *Group) Wait() {
	g.WaitGroup.Wait()
	if g.sem != nil {
		close(g.sem)
	}
}

func (g *Group) Errors() []error {
	return g.errors
}
