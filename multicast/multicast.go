package multicast

import (
	"encoding/gob"
	"fmt"
	"net"

	"golang.org/x/xerrors"
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

// multicastInterfaces はマルチキャスト可能で IPv4 アドレスを持つ UP な
// インターフェースを返す。Windows では仮想アダプタ(WSL/Hyper-V/VPN)が
// 混在しデフォルトインターフェースが当てにならないため、
// 送信・受信とも全インターフェースを対象にする
func multicastInterfaces() ([]net.Interface, error) {
	all, err := net.Interfaces()
	if err != nil {
		return nil, xerrors.Errorf("interfaces: %w", err)
	}
	res := make([]net.Interface, 0, len(all))
	for _, ifi := range all {
		if ifi.Flags&net.FlagUp == 0 || ifi.Flags&net.FlagMulticast == 0 {
			continue
		}
		addrs, err := ifi.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			if ipn, ok := a.(*net.IPNet); ok && ipn.IP.To4() != nil {
				res = append(res, ifi)
				break
			}
		}
	}
	return res, nil
}
