package udpbroadcast

import (
	"fmt"
	"net"

	"github.com/vaiojarsad/ringy/internal/appcontext"
	"github.com/vaiojarsad/ringy/internal/executor"
	netutil "github.com/vaiojarsad/ringy/internal/util/net"
)

type exec struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	c    *appcontext.AppContext
}

func NewUDPBroadcastExecutor() (executor.Executor, error) {
	c := appcontext.GetInstance()
	nc := c.CfgManager.GetNetworkConfig()
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", nc.MulticastAddr, nc.Port))
	if err != nil {
		return nil, err
	}
	iFace, err := netutil.GetDefaultInterface()
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenMulticastUDP("udp", iFace, addr)
	if err != nil {
		return nil, err
	}
	err = conn.SetReadBuffer(256)
	if err != nil {
		return nil, err
	}
	c.OutLogger.Println("Listening for multicast messages on", fmt.Sprintf("%s:%s",
		nc.MulticastAddr, nc.Port))
	return &exec{addr: addr, conn: conn, c: c}, nil
}

func (e *exec) Close() error {
	if e.conn != nil {
		return e.conn.Close()
	}
	return nil
}

func (e *exec) Do() error {
	buf := make([]byte, 256)
	n, src, err := e.conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	e.c.OutLogger.Printf("Received from %s: %s\n", src, string(buf[:n]))
	return nil
}
