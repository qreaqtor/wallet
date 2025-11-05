package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type apiServer struct {
	router *mux.Router
	addr   string

	log Log
}

func New(log Log, port int64) *apiServer {
	return &apiServer{
		router: mux.NewRouter().PathPrefix("/api").Subrouter(),
		addr:   fmt.Sprintf(":%d", port),
		log:    log,
	}
}

func (s *apiServer) Handle(method, path string, f http.HandlerFunc) {
	s.router.HandleFunc(path, f).Methods(method)
}

func (s *apiServer) Run(ctx context.Context) error {
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowCredentials(),
	)(s.router)

	server := &http.Server{
		Handler: corsHandler,
		Addr:    s.addr,
	}

	errChan := listenAndServe(server)

	s.log.Debug(ctx, "server started")
	defer s.log.Debug(ctx, "server closed")

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return server.Shutdown(shutdownCtx)
	case err := <-errChan:
		return err
	}
}

func listenAndServe(server *http.Server) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	return errChan
}
