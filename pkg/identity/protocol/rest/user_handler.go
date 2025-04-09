package rest

import (
	"money-transfer-demo/pkg/identity/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewUserHandler(r *gin.Engine) {
	groupUser := r.Group("/api/admin/user")

	groupUser.POST("/", s.createUser)
	groupUser.POST("/login", s.loginUser)
}

func (s *Server) createUser(c *gin.Context) {
	var cmd user.CreateUserCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.Dependencies.UserSvc.CreateUser(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (s *Server) loginUser(c *gin.Context) {
	var cmd user.LoginUserCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := s.Dependencies.UserSvc.LoginUser(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
