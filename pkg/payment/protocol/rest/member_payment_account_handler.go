package rest

import (
	"money-transfer-demo/pkg/identity/token"
	"money-transfer-demo/pkg/middleware"
	"money-transfer-demo/pkg/payment/memberpayacc"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) NewMemberPaymentAccountHandler(r *gin.Engine) {
	groupMem := r.Group("/api/mem/payment-account")

	groupMem.POST("/", middleware.AuthMiddleware(token.Member), s.createPaymentAccount)
	groupMem.GET("/", middleware.AuthMiddleware(token.Member), s.searchPaymentAccount)
}

func (s *Server) createPaymentAccount(ctx *gin.Context) {
	var cmd memberpayacc.CreateMemberPayAccountCommand

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := cmd.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.Dependencies.MemberPayAccSrv.Create(ctx.Request.Context(), &cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) searchPaymentAccount(ctx *gin.Context) {
	var query memberpayacc.SearchMemberPayAccountQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := s.Dependencies.MemberPayAccSrv.Search(ctx.Request.Context(), &query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
