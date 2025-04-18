package rest

import (
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/middleware"
	"money-transfer-demo/pkg/payment/withdrawal"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewWithdrawalHandler(r *gin.Engine) {
	groupMem := r.Group("/api/mem/withdrawal")
	groupAdmin := r.Group("/api/admin/withdrawal")

	groupMem.POST("/", middleware.AuthMiddleware(token.Member), s.createWithdrawal)

	groupAdmin.GET("/", middleware.AuthMiddleware(token.User), s.searchWithdrawal)
	groupAdmin.GET("/detail/:id", middleware.AuthMiddleware(token.User), s.getWithdrawalByID)
	groupAdmin.PUT("/action/:id/approve", middleware.AuthMiddleware(token.User), s.approveWithdrawal)
	groupAdmin.PUT("/action/:id/reject", middleware.AuthMiddleware(token.User), s.rejectWithdrawal)
	groupAdmin.PUT("/action/:id/transfer", middleware.AuthMiddleware(token.User), s.transferWithdrawal)
}

func (s *Server) createWithdrawal(c *gin.Context) {
	var cmd withdrawal.CreateWithdrawalCommand

	err := c.ShouldBindJSON(&cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Dependencies.WithdrawalSrv.CreateWithdrawal(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "withdrawal created")
}

func (s *Server) searchWithdrawal(c *gin.Context) {
	var query withdrawal.SearchWithdrawalQuery

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.Dependencies.WithdrawalSrv.SearchWithdrawal(c.Request.Context(), &query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) getWithdrawalByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.Dependencies.WithdrawalSrv.GetWithdrawalByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) approveWithdrawal(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Dependencies.WithdrawalSrv.ApproveWithdrawal(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) rejectWithdrawal(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Dependencies.WithdrawalSrv.RejectWithdrawal(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) transferWithdrawal(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Dependencies.WithdrawalSrv.TransferWithdrawal(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
