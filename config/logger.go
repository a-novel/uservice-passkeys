package config

import (
	_ "embed"
	"os"

	"github.com/rs/zerolog"

	"github.com/a-novel/golib/deploy"
	"github.com/a-novel/golib/loggers"
	"github.com/a-novel/golib/loggers/formatters"
)

//go:embed logger-dev.yaml
var loggerFileDev []byte

//go:embed logger-prod.yaml
var loggerFileProd []byte

type LoggerType struct {
	Dynamic bool   `yaml:"dynamic"`
	Type    string `yaml:"type"`

	Formatter formatters.Formatter
}

var Logger = deploy.LoadConfig[LoggerType](
	deploy.DevConfig(loggerFileDev),
	deploy.ProdConfig(loggerFileProd),
)

func init() {
	switch Logger.Type {
	case "console":
		logger := loggers.NewSTDOut()
		Logger.Formatter = formatters.NewConsoleFormatter(logger, !Logger.Dynamic)
	case "json":
		logger := loggers.NewZeroLog(zerolog.New(os.Stdout))
		Logger.Formatter = formatters.NewJSONFormatter(logger)
	default:
		panic("unknown logger type " + Logger.Type)
	}
}
