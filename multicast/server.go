package multicast

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"time"

	"golang.org/x/net/ipv4"
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

	group := &net.UDPAddr{IP: net.ParseIP(MulticastAddress), Port: s.UDPPort}
	log.Printf("Announce %s[%s]", s.Name, group)

	// net.Dial だと OS が選んだ1つのインターフェースからしか送信されず、
	// 仮想アダプタが混在する環境で届かないことがあるため、
	// 全インターフェースへ明示的に送信する
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return xerrors.Errorf("UDP listen: %w", err)
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)
	// 同一ホスト上のクライアントにも届くようにループバックを有効化
	if err := p.SetMulticastLoopback(true); err != nil {
		log.Printf("set multicast loopback: %v", err)
	}

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
	for {
		if err := announce(p, group, data); err != nil {
			log.Printf("announce: %+v", err)
		}
		time.Sleep(d)
	}
}

// announce はマルチキャスト可能な全インターフェースへ告知を1回ずつ送る
func announce(p *ipv4.PacketConn, group *net.UDPAddr, data []byte) error {
	ifs, err := multicastInterfaces()
	if err != nil {
		return err
	}
	sent := 0
	for i := range ifs {
		if err := p.SetMulticastInterface(&ifs[i]); err != nil {
			continue
		}
		if _, err := p.WriteTo(data, nil, group); err != nil {
			continue
		}
		sent++
	}
	if sent == 0 {
		return xerrors.New("no interface could send multicast")
	}
	return nil
}
