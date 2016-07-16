package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	defaultEndpoint       = "0.0.0.0:8000"
	defaultConfigFileName = "./dashd.json"
)

// Command line parameters
var (
	Endpoint       = defaultEndpoint
	ConfigFileName = defaultConfigFileName
	LogFileName    = ""
)

var logger *log.Logger
var logWriter io.Writer = os.Stderr

func init() {
	flag.StringVar(&Endpoint, "addr", defaultEndpoint, "HTTP endpoint")
	flag.StringVar(&LogFileName, "log", "", "Log file path")
	flag.StringVar(&ConfigFileName, "config", defaultConfigFileName, "Config file path")

	flag.Parse()

	if LogFileName != "" {
		var err error
		logWriter, err = os.Create(LogFileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while opening '%s': %s'.\n", LogFileName, err)
			os.Exit(-1)
		}
	}

	logger = log.New(logWriter, "dashd", log.Ltime|log.Lshortfile)
}
