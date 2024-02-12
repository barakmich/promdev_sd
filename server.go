package promdev_sd

import (
	"net/http"
	"sync"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	sync.Mutex
	router           *chi.Mux
	labelSets        map[string]LabelSet
	namespaceEntries map[string][]*Target
}

func NewServer() (*Server, error) {
	return &Server{
		router: chi.NewRouter(),
	}, nil
}

func (s *Server) ListenAndServe(hostport string) error {
	s.router.Use(middleware.Logger)
	s.router.Put("/register/{namespace}", s.register)
	s.router.Put("/heartbeat/{namespace}/{token}", s.heartbeat)
	s.router.Get("/discovery/{namespace}", s.discovery)
	return http.ListenAndServe(hostport, s.router)
}
