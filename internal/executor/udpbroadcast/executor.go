package udpbroadcast

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"net"
	"os"
	"path"
	"time"

	"github.com/vaiojarsad/ringy/internal/appcontext"
	"github.com/vaiojarsad/ringy/internal/executor"
	ioutil "github.com/vaiojarsad/ringy/internal/util/io"
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
	buf := make([]byte, 512)
	n, src, err := e.conn.ReadFromUDP(buf)
	if err != nil {
		return err
	}
	e.c.OutLogger.Printf("Received %d bytes from %s: %s\n", n, src, string(buf[:n]))
	e.playMP3()
	return nil
}

func (e *exec) playMP3() {
	apc := e.c.CfgManager.GetAudioPlaybackConfig()
	f, err := os.Open(path.Clean(path.Join(apc.AudioFilesPath, apc.RingAudio)))
	if err != nil {
		e.c.ErrLogger.Printf("Error opening MP3 file: %v\n", err)
		return
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		ioutil.Close(f, "playMP3")
		e.c.ErrLogger.Printf("Error decoding MP3: %v\n", err)
		return
	}
	defer ioutil.Close(streamer, "playMP3")

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		e.c.ErrLogger.Printf("Error initializing speaker: %v\n", err)
		return
	}

	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	<-done
}
