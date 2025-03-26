package rest

import (
	"money-transfer-demo/pkg/identity/member"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewMemberHandler(r *gin.Engine) {
	groupMember := r.Group("/api/mem/member")

	groupMember.POST("/", s.createMember)

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
