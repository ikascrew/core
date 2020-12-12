package multicast

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"time"

	"golang.org/x/xerrors"
)

type Client struct {
	Port     int
	Duration int
}

type ClientOption func(*Client) error

func defaultClient() *Client {
	c := Client{}
	c.Port = DefaultUDPPort
	c.Duration = 5
	return &c
}

func ClientPort(p int) func(*Client) error {
	return func(c *Client) error {
		c.Port = p
		return nil
	}
}

func ClientDuration(d int) func(*Client) error {
	return func(c *Client) error {
		c.Duration = d
		return nil
	}
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

	log.Println(add)
	addr, err := net.ResolveUDPAddr("udp", add)
	if err != nil {
		return nil, xerrors.Errorf("resolve udp address[%s]: %w", add, err)
	}

	var in *net.Interface = nil

	/*
		in, err = net.InterfaceByName("イーサネット 2")
		if err != nil {
			return nil, xerrors.Errorf("not found[%s]: %w", add, err)
		}

			interfaces, err := net.Interfaces()
			for _, elm := range interfaces {
				log.Println("Interface :" + elm.Name)
				log.Printf("Index :%d\n", elm.Index)
				log.Printf("MTU :%d\n", elm.MTU)
			}
	*/

	conn, err := net.ListenMulticastUDP("udp", in, addr)
	if err != nil {
		return nil, xerrors.Errorf("listen multicast [%s]: %w", add, err)
	}
	defer conn.Close()

	acs := make([]*AccessInfo, 0, 10)

	start := time.Now()
	dead := start.Add(time.Second * time.Duration(c.Duration+1))
	err = conn.SetDeadline(dead)

	if err != nil {
		return nil, xerrors.Errorf("set deadline: %w", err)
	}

	for {

		buffer := make([]byte, MTU)

		length, remoteAddress, err := conn.ReadFromUDP(buffer)
		if err != nil {

			var opErr *net.OpError
			if errors.As(err, &opErr) {
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
