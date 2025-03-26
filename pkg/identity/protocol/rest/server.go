package rest

import (
	"context"
	"fmt"
	"log"
	"money-transfer-demo/pkg/identity/config"
	"money-transfer-demo/pkg/identity/member"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

type Server struct {
	Cfg          *config.Config
	Dependencies *Dependencies
	Router       *gin.Engine
}

type Dependencies struct {
	MemberSvc member.Service
	Cfg       *config.Config
}

func NewServer(deps *Dependencies, cfg *config.Config) *Server {
	router := gin.Default()

	server := &Server{
		Cfg:          cfg,
		Dependencies: deps,
		Router:       router,
	}

	server.registerRoutes()

	return server
}

func (s *Server) registerRoutes() {
	r := s.Router
	s.NewMemberHandler(r)
}

func (s *Server) Run(ctx context.Context) error {
	stopCh := ctx.Done()

	router := gin.New()

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	})

	handler := c.Handler(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Cfg.Server.HTTPPort),
		Handler: handler,
	}

	go func() {
		<-stopCh
		log.Println("Shutting down HTTP server...")

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("❌ Server forced to shutdown: %v\n", err)
		}

		log.Println("✅ Server exited properly")
	}()

	log.Printf("🚀 Starting HTTP server on port %s...\n", s.Cfg.Server.HTTPPort)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
