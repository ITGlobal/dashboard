package main

// Common imports
import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	dash "github.com/itglobal/dashboard/api"
)

// Build properties
var (
	Version            = ""
	BuildConfiguration = ""
)

// Command line args
const (
	DefaultDaemonURL   = "http://127.0.0.1:8000/"
	DefaultThemeColors = "dark"
	DefaultThemeStyle  = "single"
)

var (
	DaemonURL   string
	ThemeColors string
	ThemeStyle  string
)

func init() {
	flag.StringVar(&DaemonURL, "url", DefaultDaemonURL, "dashd's URL")
	flag.StringVar(&ThemeColors, "colors", DefaultThemeColors, "Color scheme - 'dark' or 'light'")
	flag.StringVar(&ThemeStyle, "style", DefaultThemeStyle, "Visual style - 'single' or 'double'")

	flag.Parse()
}

func main() {
	flag.Parse()

	fmt.Printf("dashd by IT Global LLC %s %s\n", Version, BuildConfiguration)

	connection := dash.NewConnection(log.New(ioutil.Discard, "", 0), DaemonURL, Callback)
	defer connection.Close()

	// Configure and run UI
	theme := ThemeDefault
	if ThemeColors == "dark" {
		theme |= ThemeDark
	}
	if ThemeColors == "light" {
		theme |= ThemeLight
	}

	if ThemeStyle == "single" {
		theme |= ThemeSingle
	}
	if ThemeStyle == "double" {
		theme |= ThemeDouble
	}

	SetTheme(theme)
	Run()
}
