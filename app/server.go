package app

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/itglobal/dashboard/tile"
	log "github.com/kpango/glg"
)

type ServerParameters struct {
	Config   *Config
	Endpoint string
	WwwRoot  string
}

func NewServerParameters(config *Config) *ServerParameters {
	return &ServerParameters{
		Config:   config,
		Endpoint: "0.0.0.0:8000",
		WwwRoot:  "",
	}
}

type Server interface {
	Run() error
}

type server struct {
	parameters *ServerParameters
	manager    *tileManager
	router     *mux.Router
	providers  []tile.Provider
}

func NewServer(parameters *ServerParameters) (Server, error) {
	server := &server{
		manager:    newTileManager(),
		router:     mux.NewRouter(),
		parameters: parameters,
		providers:  make([]tile.Provider, 0),
	}

	server.initializeHTTP()
	server.initializeProviders()

	return server, nil
}

func (s *server) initializeHTTP() {
	log.Debug("Configuring HTTP API")

	s.router.Handle("/data.json", withLogging(&restAPIHandler{s.manager})).Methods("GET")
	s.router.Handle("/data.ws", &wsAPIHandler{s.manager})

	if s.parameters.WwwRoot != "" {
		handler := http.FileServer(http.Dir(s.parameters.WwwRoot))
		s.router.PathPrefix("/").Handler(withLogging(handler))
	} else {
		log.Warn("No wwwroot dir is set")
	}
}

func (s *server) initializeProviders() {
	log.Debug("Configuring tile providers")

	for _, config := range s.parameters.Config.listProviders() {
		factory := tile.GetFactory(config.Type)
		if factory == nil {
			log.Errorf("\"providers[%d]\":\tUnable to find tile provider \"%s\"", config.index, config.Type)
			continue
		}

		provider, err := factory.Create(config, s.manager)
		if err != nil {
			log.Errorf("\"providers[%d]\":\tUnable to create tile provider \"%s\"", config.index, config.Type)
			continue
		}

		err = provider.Init()
		if err != nil {
			log.Errorf("\"providers[%d]\":\tUnable to initialize tile provider \"%s\"", config.index, config.Type)
			continue
		}

		log.Infof("\"providers[%d]\":\tCreated tile provider \"%s\"", config.index, config.Type)
		s.providers = append(s.providers, provider)
	}
}

func (s *server) Run() error {
	log.Infof("Listening at \"http://%s\"", s.parameters.Endpoint)
	log.Info("Press Ctrl+C to exit")
	return http.ListenAndServe(s.parameters.Endpoint, s.router)
}
