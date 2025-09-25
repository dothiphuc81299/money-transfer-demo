package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewHealthHandler(r *gin.Engine) {
	r.GET("/api/health", s.health)
}

func (s *Server) health(c *gin.Context) {
	status := "ok"
	if err := s.Dependencies.Postgres.Ping(); err != nil {
		status = "unhealthy"
	}

	// response phẳng, dễ parse bằng jq
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
