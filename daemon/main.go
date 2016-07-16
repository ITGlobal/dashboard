package main

// Common imports
import (
	"log"
	"net/http"

	_ "expvar"
)

// Import providers
import (
	_ "github.com/itglobal/dashboard/providers/mongodb"
	_ "github.com/itglobal/dashboard/providers/ping"
	_ "github.com/itglobal/dashboard/providers/sim"
	_ "github.com/itglobal/dashboard/providers/teamcity"
)

// Build properties
var (
	Version            = ""
	BuildConfiguration = ""
)

func main() {
	logger.Printf("dashd by IT Global LLC %s %s\n", Version, BuildConfiguration)

	// Create providers
	Providers = createProviders(updateData)

	// Initialize HTTP API
	r := http.NewServeMux()
	r.HandleFunc("/", GetDataHandler)
	h := CreateHandler(log.New(logWriter, "http", log.Ltime|log.Lshortfile), r)

	// Run
	logger.Printf("dashd is running at endpoint %s", Endpoint)
	logger.Fatal(http.ListenAndServe(Endpoint, h))
}
