package logger

import (
	"os"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New (level string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).Level(lvl).With().Timestamp().Caller().Logger()
	return logger
}

func SetupGlobal(level string) {
	log.Logger = New(level)
}