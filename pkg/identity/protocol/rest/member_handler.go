package rest

import (
	"money-transfer-demo/pkg/identity/member"
	"money-transfer-demo/pkg/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewMemberHandler(r *gin.Engine) {
	groupMember := r.Group("/api/mem/member")

	groupMember.POST("/", s.createMember)
	groupMember.POST("/login", s.loginMember)
	groupMember.GET("/detail/:id", s.getMemberByID, middleware.AuthMiddleware())

}

func (h *Server) createMember(c *gin.Context) {
	var cmd member.CreateMemberCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.Dependencies.MemberSvc.CreateMember(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func (h *Server) loginMember(c *gin.Context) {
	var cmd member.LoginMemberCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := h.Dependencies.MemberSvc.LoginMember(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

func (h *Server) getMemberByID(c *gin.Context) {
	id := c.Param("id")

	member, err := h.Dependencies.MemberSvc.GetMemberByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}
