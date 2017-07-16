package main

import (
	"github.com/itglobal/dashboard/app"
	log "github.com/kpango/glg"
	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Config string `cli:"c,config" dft:"$DASH_CONFIG"`
	Addr   string `cli:"addr" dft:"$DASH_ENDPOINT"`
}

const (
	defaultEndpoint       = "0.0.0.0:8000"
	defaultConfigFileName = "./dashboard.json"
)

var (
	endpoint       = defaultEndpoint
	configFileName = defaultConfigFileName
)

func main() {
	cli.SetUsageStyle(cli.ManualStyle)
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		log.Info("Dashboard by IT Global LLC")

		endpoint := argv.Addr
		if endpoint == "" {
			endpoint = defaultEndpoint
		}

		configPath := argv.Config
		if configPath == "" {
			configPath = defaultConfigFileName
		}

		config, err := app.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Unable to configure application. %s", err)
			return err
		}

		parameters := app.NewServerParameters(config)
		parameters.Endpoint = endpoint
		parameters.WwwRoot = "./www"

		server, err := app.NewServer(parameters)
		if err != nil {
			log.Fatalf("Unable to create server. %s", err)
			return err
		}

		err = server.Run()
		if err != nil {
			log.Fatalf("Critical error. %s", err)
			return err
		}

		return err
	})
}
