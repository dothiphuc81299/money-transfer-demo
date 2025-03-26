package rest

import (
	"context"
	"fmt"
	"log"
	"money-transfer-demo/pkg/payment/config"
	"money-transfer-demo/pkg/payment/memberacc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

type Server struct {
	Cfg          *config.Config
	Dependencies *Dependencies
	Router       *gin.Engine
	HTTPServer   *http.Server
}

type Dependencies struct {
	MemberAccSvc memberacc.Service
	Cfg          *config.Config
}

func NewServer(deps *Dependencies, cfg *config.Config) *Server {
	router := gin.Default()

	server := &Server{
		Cfg:          cfg,
		Dependencies: deps,
		Router:       router,
	}

	return server
}

func (s *Server) registerRoutes(router *gin.Engine) {

}

func (s *Server) Run(ctx context.Context) error {
	stopCh := ctx.Done()

	router := gin.New()

	s.registerRoutes(router)
	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	handler := c.Handler(router)

	s.HTTPServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Cfg.Server.HTTPPort),
		Handler: handler,
	}

	go func() {
		<-stopCh
		log.Println("Shutting down HTTP server...")

		if err := s.HTTPServer.Shutdown(context.Background()); err != nil {
			log.Printf("❌ Server forced to shutdown: %v\n", err)
		}

		log.Println("✅ Server exited properly")
	}()

	log.Printf("🚀 Starting HTTP server on port %s...\n", s.Cfg.Server.HTTPPort)

	if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.HTTPServer != nil {
		return s.HTTPServer.Shutdown(ctx)
	}
	return nil
}
