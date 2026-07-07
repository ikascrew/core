package multicast

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/ipv4"
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
	addr, err := net.ResolveUDPAddr("udp4", add)
	if err != nil {
		return nil, xerrors.Errorf("resolve udp address[%s]: %w", add, err)
	}

	// ListenMulticastUDP(nil) はデフォルトインターフェースにしか
	// join しないため(Windows の複数NIC環境で受信できない原因)、
	// マルチキャスト可能な全インターフェースへ追加で join する
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, xerrors.Errorf("listen multicast [%s]: %w", add, err)
	}
	defer conn.Close()

	p := ipv4.NewPacketConn(conn)
	group := &net.UDPAddr{IP: addr.IP}
	if ifs, err := multicastInterfaces(); err == nil {
		for i := range ifs {
			// デフォルトインターフェースは join 済みでエラーになるが問題ない
			p.JoinGroup(&ifs[i], group)
		}
	}
	// 同一ホストのサーバーも発見できるようにループバックを有効化
	if err := p.SetMulticastLoopback(true); err != nil {
		log.Printf("set multicast loopback: %v", err)
	}

	acs := make([]*AccessInfo, 0, 10)
	seen := make(map[string]bool)

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

		// 全インターフェースへ告知しているため同じサーバーから
		// 複数回届く。同一の告知は捨てる
		key := fmt.Sprintf("%s|%d|%s|%d", m.Name, m.Type, m.Address, m.Port)
		if seen[key] {
			continue
		}
		seen[key] = true

		acs = append(acs, &m)
	}
}
