package multicast

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"

	"golang.org/x/xerrors"
)

type Server struct {
	Name      string
	Type      ServerType
	UDPPort   int
	UsingPort int
	Duration  int
}

type ServerOption func(*Server) error

func ServerName(name string) ServerOption {
	return func(s *Server) error {
		s.Name = name
		return nil
	}
}

func Port(p int) ServerOption {
	return func(s *Server) error {
		s.UDPPort = p
		return nil
	}
}

func Use(p int) ServerOption {
	return func(s *Server) error {
		s.UsingPort = p
		return nil
	}
}

func Duration(d int) ServerOption {
	return func(s *Server) error {
		s.Duration = d
		return nil
	}
}

func Type(t ServerType) ServerOption {
	return func(s *Server) error {
		s.Type = t
		if s.UsingPort == DefaultServerPort {
			switch t {
			case TypeIkasbox:
				s.UsingPort = DefaultDatabasePort
			}
		}
		return nil
	}
}

func defaultServer() *Server {
	s := Server{}
	s.Name = "ikascrew"
	s.UDPPort = DefaultUDPPort
	s.Type = TypeIkascrew
	s.UsingPort = DefaultServerPort
	s.Duration = 5
	return &s
}

func NewServer(opts ...ServerOption) (*Server, error) {
	s := defaultServer()
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			log.Println(err)
		}
	}
	return s, nil
}

func (s *Server) Dial() error {

	ser := createAddress(s.UDPPort)
	log.Printf("Dial %s[%s]", s.Name, ser)

	conn, err := net.Dial("udp", ser)
	if err != nil {
		return xerrors.Errorf("UDP Dial: %w", err)
	}
	defer conn.Close()

	msg := AccessInfo{}
	msg.Name = s.Name
	msg.Port = s.UsingPort
	msg.Address = ""
	msg.Type = s.Type

	var w bytes.Buffer
	enc := gob.NewEncoder(&w)
	err = enc.Encode(msg)
	if err != nil {
		return xerrors.Errorf("encode : %w", err)
	}

	data := w.Bytes()
	if len(data) > MTU {
		return xerrors.Errorf("MTU Size Error: %d", MTU)
	}

	d := time.Duration(s.Duration) * time.Second
	ticker := time.Tick(d)
	for range ticker {
		conn.Write(data)
	}

	return nil
}
