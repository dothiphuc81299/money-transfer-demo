package rest

import (
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/middleware"
	"money-transfer-demo/pkg/payment/deposit"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewDepositHandler(r *gin.Engine) {
	groupMem := r.Group("/api/mem/deposit")
	groupAdmin := r.Group("/api/admin/deposit")

	groupMem.POST("/lbt", middleware.AuthMiddleware(token.Member), s.createDeposit)

	groupAdmin.PUT("/action/:depositId/lbt/approve", middleware.AuthMiddleware(token.User), s.approveLBTDeposit)
	groupAdmin.PUT("/action/:depositId/lbt/reject", middleware.AuthMiddleware(token.User), s.rejectLBT)
}

func (s *Server) createDeposit(c *gin.Context) {
	var cmd deposit.CreateDepositCommand

	memberID, exists := c.Get("member_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	cmd.MemberID = memberID.(int64)

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.Dependencies.DepositSrv.CreateDeposit(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) approveLBTDeposit(c *gin.Context) {
	var cmd deposit.UpdateDepositStatusCommand

	depositID, err := strconv.ParseInt(c.Param("depositId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ID = depositID
	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Dependencies.DepositSrv.ApproveLBT(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) rejectLBT(c *gin.Context) {
	var cmd deposit.UpdateDepositStatusCommand

	depositID, err := strconv.ParseInt(c.Param("depositId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.ID = depositID
	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Dependencies.DepositSrv.RejectLBT(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
