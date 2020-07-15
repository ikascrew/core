package multicast

import (
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(AccessInfo{})
}

const (
	MulticastAddress    = "224.0.0.224"
	DefaultUDPPort      = 15496
	DefaultServerPort   = 55555
	DefaultDatabasePort = 5555
	MTU                 = 1500
)

type ServerType int

const (
	TypeIkascrew ServerType = iota
	TypeIkasbox
)

func (t ServerType) Name() string {
	switch t {
	case TypeIkascrew:
		return "ikascrew server"
	case TypeIkasbox:
		return "ikasbox"
	}
	return "Type Not Found"
}

func createAddress(p int) string {
	return fmt.Sprintf("%s:%d", MulticastAddress, p)
}

type AccessInfo struct {
	Name    string
	Type    ServerType
	Address string
	Port    int
}

func (m AccessInfo) String() string {
	return fmt.Sprintf("%s(%s)=[%s:%d]", m.Name, m.Type.Name(), m.Address, m.Port)
}
