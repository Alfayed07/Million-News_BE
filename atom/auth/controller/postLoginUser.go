package atom_auth

import (
	atom "BACKEND_SEJUTA_BERITA/atom/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PostLoginUser handles POST /auth/login
func PostLoginUser(ctx *gin.Context) {
	var req atom.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "invalid request body",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "username/password cannot be empty",
		})
		return
	}

	user, token, err := atom.LoginUseCase(req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  401,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "successfully logged in",
		"token":   token,
		"user":    user,
	})
}
