package main

// Common imports
import (
	"flag"
	"io"
	"log"
	"os"

	dash "github.com/itglobal/dashboard/api"

	_ "expvar"
	"net/http"
)

// Import providers
import (
	_ "github.com/itglobal/dashboard/providers/mongodb"
	_ "github.com/itglobal/dashboard/providers/ping"
	_ "github.com/itglobal/dashboard/providers/sim"
	_ "github.com/itglobal/dashboard/providers/teamcity"
)

var Version = ""
var BuildConfiguration = ""

var logFileName = defaultLogFileName

const defaultLogFileName = "dashboard.log"

type splitWriter struct {
	enableTerminal bool
	file           io.Writer
}

var logWriter splitWriter

func (w *splitWriter) Write(p []byte) (n int, err error) {
	if w.enableTerminal {
		os.Stdout.Write(p)
	}

	return w.file.Write(p)
}

func init() {
	flag.StringVar(&configFileName, "config", defaultConfigFileName, "Config file path")
	flag.StringVar(&logFileName, "log", defaultLogFileName, "Log file path")
}

func setupLog() {
	file, err := os.Create(logFileName)
	if err != nil {
		panic(err)
	}

	logWriter = splitWriter{true, file}
	log.SetOutput(&logWriter)
}

func disableConsoleLog() {
	logWriter.enableTerminal = false
}

func main() {
	flag.Parse()
	go func() {
		http.ListenAndServe(":8080", nil)
	}()

	setupLog()

	log.Printf("[dash] Dashboard by IT Global LLC %s %s\n", Version, BuildConfiguration)

	// Read config
	log.Printf("[dash] reading config file '%s'", configFileName)
	config, err := readConfig()
	if err != nil {
		log.Printf("[dash] Unable to read config file '%s'. %s", configFileName, err)
		os.Exit(-1)
		return
	}

	// Create providers
	var providers []dash.Provider
	for i, providerJSON := range config.Providers {
		providerConf, err := newConfigReader(providerJSON)
		if err != nil {
			log.Printf("[dash] Error in 'providers[%d]': %s", i, err)
			continue
		}

		t, err := providerConf.GetString("type")
		if err != nil {
			log.Printf("[dash] Error in 'providers[%d]': property 'type' is missing", i)
			continue
		}

		factory := dash.GetFactory(t)
		if factory == nil {
			log.Printf("[dash] Error in 'providers[%d]': no such provider - '%s'", i, t)
			continue
		}

		provider, err := factory(providerConf, Callback)
		if err != nil {
			log.Printf("[dash] Error in 'providers[%d]': unable to create a provider. %s", i, err)
			continue
		}

		providers = append(providers, provider)
		log.Printf("[dash] Added provider '%s'", t)
	}

	// Configure and run UI
	log.Printf("[dash] Starting UI")
	theme := ThemeDefault
	if config.Theme.Colors == "dark" {
		theme |= ThemeDark
	}

	if config.Theme.Colors == "light" {
		theme |= ThemeLight
	}

	if config.Theme.Style == "single" {
		theme |= ThemeSingle
	}

	if config.Theme.Style == "double" {
		theme |= ThemeDouble
	}

	SetTheme(theme)

	// Shuw down terminal logging
	disableConsoleLog()

	Run()
}
