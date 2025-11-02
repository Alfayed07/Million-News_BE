package atom_auth

import (
	auth "BACKEND_SEJUTA_BERITA/atom/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostResetPassword(ctx *gin.Context){
	var req auth.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message":"invalid request body"}); return
	}
	if err := auth.ResetPasswordUseCase(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return
	}
	ctx.JSON(http.StatusOK, gin.H{"message":"password has been reset"})
}
