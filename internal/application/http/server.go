package http

import (
	_ "abt-dashboard-api/docs"
	"abt-dashboard-api/internal/application/http/getCountryRevenue"
	"abt-dashboard-api/internal/application/http/getMonthlySalesVolume"
	"abt-dashboard-api/internal/application/http/getTopProducts"
	"abt-dashboard-api/internal/application/http/getTopRegions"
	"abt-dashboard-api/internal/application/http/ping"
	pkgErrors "abt-dashboard-api/pkg/errors"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/cors"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Server *http.Server
	DbConn *sql.DB
}

func NewServer(abtDashboardDBConn *sql.DB) *Server {
	return &Server{
		DbConn: abtDashboardDBConn,
	}
}

func (s *Server) Start(ctx context.Context) {

	r := s.registerRouter()

	port := 8080
	readTimeout := 5 * time.Second
	writeTimeout := 10 * time.Second
	idleTimeout := 15 * time.Second

	s.Server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	go func() {
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %d", port)

}

func (s *Server) Stop(ctx context.Context) {
	err := s.Server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("transport.server.Stop.error stopping transport server : %v", err)
	}
	log.Println("Server gracefully stopped")
}

func (s *Server) registerRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(SimpleContextMiddleware)
	r.Use(LoggingMiddleware)
	r.Use(RecoveryMiddleware)

	r.Get("/ping", ping.Handler)
	r.Get("/v1/metrics/country-revenue", getCountryRevenue.Handler(s.DbConn))
	r.Get("/v1/metrics/top-products", getTopProducts.Handler(s.DbConn))
	r.Get("/v1/metrics/monthly-sales", getMonthlySalesVolume.Handler(s.DbConn))
	r.Get("/v1/metrics/top-regions", getTopRegions.Handler(s.DbConn))

	// Swagger Docs
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		pkgErrors.EncodeError(w, pkgErrors.New("route not found", 404))
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		pkgErrors.EncodeError(w, pkgErrors.New("method not allowed", 405))
	})

	return r
}

func SimpleContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		log.Printf("%s %s %d %s %s", r.Method, r.URL.Path, rw.statusCode, duration, ip)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				pkgErrors.EncodeError(w, pkgErrors.New(fmt.Sprintf("internal error: %v", err), 500))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}
