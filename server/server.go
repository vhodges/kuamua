package server

import (
    "context"
	"encoding/json"
    "fmt"
	"net/http"
    "log"
    "os"
	"os/signal"
	"syscall"

    "github.com/vhodges/kuamua/database"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

	"github.com/jackc/pgx/v5/pgxpool"

	//"quamina.net/go/quamina"
)

type Server struct {
	pgxconnection *pgxpool.Pool // Database connection pool
	Db *database.Queries    // Db/Query interface
	Mux *chi.Mux            // http router

	store *QuaminaStore
}

func New() *Server {

	// urlExample := "postgres://Patternname:password@localhost:5432/database_name"
	conn, err := pgxpool.New(context.Background(), os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	db := database.New(conn)

	store, store_err := NewQuaminaStore(conn, 10000) // DB and How large is lru cache size
	if (store_err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to create quamina store: %v\n", err)
		os.Exit(1)
	}

	//mux := http.NewServeMux()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return &Server{pgxconnection: conn, Db: db, Mux: r, store: store}
}

func (service *Server) Port() string {
	s := os.Getenv("SERVERPORT")
	if s == "" {
	  s = "3000"
	}
	return s
}

func (service *Server) start() {

	srv := &http.Server{
		Addr:    ":" +service.Port(),
		Handler: service.Mux,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	log.Print("The service is ready to listen and serve.")

	killSignal := <-interrupt
	switch killSignal {
	case os.Kill:
		log.Print("Got SIGKILL...")
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}

	log.Print("The service is shutting down...")
	srv.Shutdown(context.Background())
	log.Print("Done")
}

func (service *Server) Run(crudenabled bool) {

	defer service.pgxconnection.Close()

	// Register our routes

	// Handler for finding matching patterns for the document/message in Body
	service.Mux.Post("/document/patterns", service.matchHandler)

	// Crud is disabled by default (eg if the caller(s) are managing the Patterns)
	if crudenabled {
		// RESTy routes for "Patterns" resource
		service.Mux.Route("/patterns", func(r chi.Router) {
			r.Get("/{owner}/{group}", service.listPatternsHandler)
			r.Post("/", service.createPatternHandler)	// POST /Patterns

			r.Route("/{PatternID}", func(r chi.Router) {
				r.Use(service.patternCtx)            			// Load the *Pattern on the request context
				r.Get("/", service.getPatternHandler)      		// GET /Patterns/123
				r.Put("/", service.updatePatternHandler)    	// PUT /Patterns/123
				r.Delete("/", service.deletePatternHandler) 	// DELETE /Patterns/123
			})
		})
	}
	// Routes for K8s endpoints
	service.Mux.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	service.Mux.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	
	service.start()
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
