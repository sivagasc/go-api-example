package common

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var logger *zerolog.Logger
var once sync.Once

// GetLoggerInstance ...
func GetLoggerInstance(filename string) *zerolog.Logger {
	once.Do(func() {
		// viper.SetConfigFile(EnvFile)

		logger = createLogger(filename)
	})
	return logger
}

func createLogger(fname string) *zerolog.Logger {

	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	// defer file.Close() // TODO: Check it

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	multi := zerolog.MultiLevelWriter(consoleWriter, file)
	logger := zerolog.New(multi).With().Timestamp().Logger()

	// if env == ProductionEnv {
	// 	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	// 	logger.Info().Msg("*** Production Configuration ***")
	// } else {
	// 	zl := logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// 	logger = &zl
	// 	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// 	logger.Info().Msg("*** Non-production Configuration ***")
	// 	logger.Debug().Msg("*** Debug Logging Enabled ***")
	// }

	return &logger
}
