package appcontext

import (
	"github.com/vaiojarsad/ringy/internal/config"
	"log"
	"sync"
)

type AppContext struct {
	WG *sync.WaitGroup // the main go routine will wait on this wg on exiting, so, all goroutines should either
	// add to this wg or be waited by one that does so
	Quit       chan bool
	OutLogger  *log.Logger
	ErrLogger  *log.Logger
	CfgManager config.Manager
}

var (
	instance *AppContext
	once     sync.Once
)

func GetInstance() *AppContext {
	once.Do(func() {
		var wg sync.WaitGroup
		instance = &AppContext{
			WG:   &wg,
			Quit: make(chan bool),
		}
	})
	return instance
}
