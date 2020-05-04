package common

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var logger *zerolog.Logger
var once sync.Once

// GetLoggerInstance ...
func GetLoggerInstance() *zerolog.Logger {
	return logger
}

// SetupLoggerInstance ...
func SetupLoggerInstance(filename, env string) *zerolog.Logger {
	once.Do(func() {
		logger = createLogger(filename, env)
	})
	return logger
}

func createLogger(fname string, env string) *zerolog.Logger {

	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	// defer file.Close() // TODO: Need to close this file

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, file)
	logger := zerolog.New(multi).With().Timestamp().Logger()

	if env == "Production" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger.Info().Msg("*** Production Configuration ***")
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logger.Info().Msg("*** Non-production Configuration ***")
		logger.Debug().Msg("*** Debug Logging Enabled ***")
	}

	return &logger
}
