package multicast

import (
	"bytes"
	"encoding/gob"
	"net"
	"time"

	"golang.org/x/xerrors"
)

type Client struct {
	Port int
}

type ClientOption func(*Client) error

func defaultClient() *Client {
	c := Client{}
	c.Port = DefaultUDPPort
	return &c
}

func NewClient(opts ...ClientOption) (*Client, error) {
	c := defaultClient()
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
		}
	}
	return c, nil
}

func (c *Client) Find() ([]*AccessInfo, error) {

	add := createAddress(c.Port)

	addr, err := net.ResolveUDPAddr("udp", add)
	if err != nil {
		return nil, xerrors.Errorf("resolve udp address[%s]: %w", add, err)
	}

	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		return nil, xerrors.Errorf("listen multicast [%s]: %w", add, err)
	}
	defer conn.Close()

	acs := make([]*AccessInfo, 0, 10)

	start := time.Now()
	dead := start.Add(time.Second * 11)
	err = conn.SetDeadline(dead)
	if err != nil {
		return nil, xerrors.Errorf("set deadline: %w", err)
	}

	for {

		buffer := make([]byte, MTU)

		length, remoteAddress, err := conn.ReadFromUDP(buffer)
		if err != nil {

			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Timeout() {
					return acs, nil
				}
				return nil, xerrors.Errorf("net OpError(not timeout): %w", err)
			} else {
				return nil, xerrors.Errorf("readFormUDP: %w", err)
			}
		}

		m := AccessInfo{}

		b := bytes.NewBuffer(buffer[:length])
		d := gob.NewDecoder(b)
		err = d.Decode(&m)
		if err != nil {
			return nil, xerrors.Errorf("decode: %w", err)
		}

		m.Address = remoteAddress.IP.String()

		acs = append(acs, &m)
	}

	return acs, nil
}
