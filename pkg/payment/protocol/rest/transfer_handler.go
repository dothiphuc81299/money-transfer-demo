package rest

import (
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/middleware"
	"money-transfer-demo/pkg/payment/transfer"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewTransferHandler(r *gin.Engine) {
	groupAdmin := r.Group("/api/admin/transfer")
	groupMem := r.Group("/api/mem/transfer")

	groupMem.POST("/", middleware.AuthMiddleware(token.Member), s.createTransfer)
	groupMem.GET("/", middleware.AuthMiddleware(token.Member), s.searchTransfer)
	groupMem.GET("/detail/:id", middleware.AuthMiddleware(token.Member), s.getTransferByID)

	groupAdmin.GET("/", middleware.AuthMiddleware(token.User), s.searchTransfer)
	groupAdmin.PUT("/detail/:id", middleware.AuthMiddleware(token.User), s.updateTransfer)
	groupAdmin.GET("/detail/:id", middleware.AuthMiddleware(token.User), s.getTransferByID)
}

func (s *Server) createTransfer(c *gin.Context) {
	var cmd transfer.CreateTransferCommand

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.Dependencies.TransferSrv.Create(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) searchTransfer(c *gin.Context) {
	var query transfer.SearchTransferQuery

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.Dependencies.TransferSrv.Search(c.Request.Context(), &query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) getTransferByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.Dependencies.TransferSrv.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) updateTransfer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cmd transfer.UpdateTransferStatusCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ID = id
	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Dependencies.TransferSrv.UpdateStatus(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
