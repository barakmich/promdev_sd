package promdev_sd

import (
	"log"
	"net/http"
	"sync"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	sync.Mutex
	router           *chi.Mux
	labelSets        map[string]LabelSet
	namespaceEntries map[string][]*Target
	config           Config
}

type Config struct {
	TimeoutPeriod time.Duration
	GCPeriod      time.Duration
}

var DefaultConfig = Config{
	TimeoutPeriod: 2 * time.Minute,
	GCPeriod:      30 * time.Second,
}

func NewServer(config Config) (*Server, error) {
	return &Server{
		router:           chi.NewRouter(),
		labelSets:        make(map[string]LabelSet),
		namespaceEntries: make(map[string][]*Target),
		config:           config,
	}, nil
}

func (s *Server) ListenAndServe(hostport string) error {
	s.router.Use(middleware.Logger)
	s.router.Put("/register/{namespace}", s.register)
	s.router.Put("/heartbeat/{namespace}/{token}", s.heartbeat)
	s.router.Delete("/heartbeat/{namespace}/{token}", s.delete)
	s.router.Get("/discovery/{namespace}", s.discovery)
	log.Println("Serving on", hostport)
	return http.ListenAndServe(hostport, s.router)
}
