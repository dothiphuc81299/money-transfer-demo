package rest

import (
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/middleware"
	"money-transfer-demo/pkg/payment/bankacc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewBankAccountHandler(r *gin.Engine) {
	groupAdmin := r.Group("/api/admin/bank-account")

	groupAdmin.POST("/", s.createBankAccount, middleware.AuthMiddleware(token.User))
}

func (s *Server) createBankAccount(c *gin.Context) {
	var cmd bankacc.CreateBankAccountCommand

	err := c.ShouldBindJSON(&cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cmd.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = s.Dependencies.BankAccSrv.CreateBankAccount(c.Request.Context(), &cmd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
