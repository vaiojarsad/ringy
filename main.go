package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vaiojarsad/ringy/internal/appcontext"
	"github.com/vaiojarsad/ringy/internal/config"
	"github.com/vaiojarsad/ringy/internal/executor"
	"github.com/vaiojarsad/ringy/internal/executor/udpbroadcast"
	ioutil "github.com/vaiojarsad/ringy/internal/util/io"
)

func main() {
	c := appcontext.GetInstance()
	cfgManager, err := config.GetInstance()
	if err != nil {
		log.Printf("error getting configuration... %v\n", err)
		os.Exit(1)
	}
	lc := cfgManager.GetLoggerConfig()
	// Set up a logger for output
	c.ErrLogger = log.New(os.Stderr, lc.Prefix, lc.Flag)
	c.OutLogger = log.New(os.Stdout, lc.Prefix, lc.Flag)

	// Set up a channel to handle OS signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Print a message to indicate the daemon has started
	c.OutLogger.Println("Starting daemon...")

	e, err := udpbroadcast.NewUDPBroadcastExecutor()
	if err != nil {
		c.ErrLogger.Printf("error creating executor... %v\n", err)
		os.Exit(1)
	}

	launch(e)

	<-signals
	close(c.Quit)
	cleanExit(e)
}

func cleanExit(e executor.Executor) {
	ioutil.Close(e, "main")
	c := appcontext.GetInstance()
	c.OutLogger.Println("Shutting down daemon...")
	c.OutLogger.Println("start waiting for the wait group to be done")

	done := make(chan bool)
	go func() {
		defer close(done)
		c.WG.Wait()
	}()

	sleepTimeInSeconds := 1
	t := time.NewTimer(time.Duration(sleepTimeInSeconds) * time.Second)
	waits := 10

	for i := 0; i < waits; i++ {
		select {
		case <-done:
			c.OutLogger.Println("wait group is done")
			return
		case <-t.C:
			t.Reset(time.Duration(sleepTimeInSeconds) * time.Second)
			c.OutLogger.Println("wait group is not done yet")
		}
	}
}

func launch(e executor.Executor) {
	c := appcontext.GetInstance()
	c.WG.Add(1)
	go func() {
		for {
			select {
			case <-c.Quit:
				c.WG.Done()
				return
			default:
				err := e.Do()
				if err != nil {
					c.OutLogger.Printf("error processing... %v", err)
				}
			}
		}
	}()
}
