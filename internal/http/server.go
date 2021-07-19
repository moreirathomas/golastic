package http

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/moreirathomas/golastic/internal/repository"
	"github.com/moreirathomas/golastic/pkg/logger"
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
	s.router.Use(logger.Middleware)
	s.registerRoutes()
}

// registerRoutes registers each entity's routes on the server.
func (s *Server) registerRoutes() {
	const bookID = "{bookID:[a-zA-Z0-9_-]+}"

	// Root
	s.router.HandleFunc("/", s.handleIndex)

	// Insert new book
	s.router.HandleFunc("/books", s.InsertBook).Methods(http.MethodPost)

	// Search books (query)
	s.router.HandleFunc("/books", s.SearchBooks).Methods(http.MethodGet)

	// Get book by ID
	s.router.HandleFunc("/books/"+bookID, s.GetBookByID).Methods(http.MethodGet)

	// Update book
	s.router.HandleFunc("/books/"+bookID, s.UpdateBook).Methods(http.MethodPut)

	// Delete book by ID
	s.router.HandleFunc("/books/"+bookID, s.DeleteBook).Methods(http.MethodDelete)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text-plain")
	w.WriteHeader(200)
	w.Write([]byte("Welcome to Golastic!"))
}
