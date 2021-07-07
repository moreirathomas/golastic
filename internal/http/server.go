package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/moreirathomas/golastic/internal/repository"
	"github.com/moreirathomas/golastic/pkg/httputil"
)

// Server represents the main server for the API.
type Server struct {
	*http.Server
	router     *mux.Router
	Repository repository.Repository
}

// NewServer returns a new Server given configuration parameters.
func NewServer(addr string, repo repository.Repository) *Server {
	return &Server{
		Server:     &http.Server{Addr: addr},
		Repository: repo,
	}
}

// Start launches the server.
// It serves its attached router at its Addr.
func (s *Server) Start() error {
	s.initRouter()
	s.Handler = s.router

	log.Printf("Server listening at http://localhost%s\n", s.Addr)

	return s.ListenAndServe()
}

func (s *Server) initRouter() {
	s.router = mux.NewRouter().StrictSlash(true)
	s.router.Use(httputil.RequestLogger)
	s.registerRoutes()
}

// registerRoutes registers each entity's routes on the server.
func (s *Server) registerRoutes() {
	s.router.HandleFunc("/", s.handleIndex)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text-plain")
	w.WriteHeader(200)
	w.Write([]byte("Welcome to Golastic!"))
}
