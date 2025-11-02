package atom_auth

import (
	auth "BACKEND_SEJUTA_BERITA/atom/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostForgotPassword(ctx *gin.Context){
	var req auth.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"invalid request body"}); return
	}
	token, err := auth.ForgotPasswordUseCase(req)
	if err != nil { ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
	// In production you would not return the token; this is for testing/demo
	ctx.JSON(http.StatusOK, gin.H{"message":"reset email sent (dev)", "reset_token": token})
}
