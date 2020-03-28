package main

import (
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
)

var l *log.Logger

// TODO grab values from the configuration.
func init() {
	l = &log.Logger{
		Handler: cli.New(os.Stdout),
		Level:   log.InfoLevel,
	}
}

func SetLogLevel(lstr string) {
	l.Level = log.MustParseLevel(lstr)
}
