package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type environmentBackedManager struct {
	loggerConfig        *LoggerConfig
	networkConfig       *NetworkConfig
	audioPlaybackConfig *AudioPlaybackConfig
}

func (m *environmentBackedManager) GetLoggerConfig() *LoggerConfig {
	return m.loggerConfig
}

func (m *environmentBackedManager) GetNetworkConfig() *NetworkConfig {
	return m.networkConfig
}

func (m *environmentBackedManager) GetAudioPlaybackConfig() *AudioPlaybackConfig {
	return m.audioPlaybackConfig
}

func newEnvironmentBackedManager() (Manager, error) {
	loggerConfig, err := getLoggerConfig()
	if err != nil {
		return nil, err
	}

	networkConfig, err := getNetworkConfig()
	if err != nil {
		return nil, err
	}

	audioPlaybackConfig, err := getAudioPlaybackConfig()
	if err != nil {
		return nil, err
	}

	m := &environmentBackedManager{
		loggerConfig:        loggerConfig,
		networkConfig:       networkConfig,
		audioPlaybackConfig: audioPlaybackConfig,
	}
	return m, nil
}

func getLoggerConfig() (*LoggerConfig, error) {
	var flag = log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Llongfile
	flagStr := os.Getenv("LOGGER_FLAG")
	flagStr = strings.TrimSpace(flagStr)
	if flagStr != "" {
		i, err := strconv.Atoi(flagStr)
		if err != nil {
			return nil, err
		}
		if i > 0 && i < 128 {
			flag = i
		}
	}

	prefix := os.Getenv("LOGGER_PREFIX")
	prefix = strings.TrimLeftFunc(prefix, unicode.IsSpace)

	return &LoggerConfig{
		Flag:   flag,
		Prefix: prefix,
	}, nil
}

func getNetworkConfig() (*NetworkConfig, error) {
	ma := os.Getenv("NET_MULTI_CAST_ADDR")
	ma = strings.TrimLeftFunc(ma, unicode.IsSpace)

	p := os.Getenv("NET_PORT")
	p = strings.TrimLeftFunc(p, unicode.IsSpace)

	return &NetworkConfig{
		MulticastAddr: ma,
		Port:          p,
	}, nil
}

func getAudioPlaybackConfig() (*AudioPlaybackConfig, error) {
	audioFilesPath := os.Getenv("AUDIO_FILES_PATH")
	audioFilesPath = strings.TrimLeftFunc(audioFilesPath, unicode.IsSpace)

	ringAudio := os.Getenv("RING_AUDIO")
	ringAudio = strings.TrimLeftFunc(ringAudio, unicode.IsSpace)

	return &AudioPlaybackConfig{
		AudioFilesPath: audioFilesPath,
		RingAudio:      ringAudio,
	}, nil
}
