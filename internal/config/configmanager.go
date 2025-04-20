package config

import (
	"sync"
)

var (
	instance Manager
	once     sync.Once
)

// Manager handle the different configurations
type Manager interface {
	GetLoggerConfig() *LoggerConfig
	GetNetworkConfig() *NetworkConfig
}

type LoggerConfig struct {
	Prefix string
	Flag   int
}

type NetworkConfig struct {
	MulticastAddr string
	Port          string
}

// GetInstance returns the singleton instance, creating it if necessary
func GetInstance() (Manager, error) {
	var err error = nil
	once.Do(func() {
		instance, err = newEnvironmentBackedManager()
	})
	return instance, err
}
