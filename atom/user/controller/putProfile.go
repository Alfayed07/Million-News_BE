package atom_user

import (
	atom "BACKEND_SEJUTA_BERITA/atom/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PutProfile(c *gin.Context) {
    userID, err := subjectFromAuth(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
        return
    }
    var req atom.UpdateProfileRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
        return
    }
    prof, err := atom.UpdateProfileUseCase(userID, req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, prof)
}
