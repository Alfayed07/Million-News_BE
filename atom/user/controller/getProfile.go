package atom_user

import (
	atom "BACKEND_SEJUTA_BERITA/atom/user"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetProfile(c *gin.Context) {
    userID, err := subjectFromAuth(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
        return
    }
    prof, err := atom.GetProfileUseCase(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, prof)
}

func subjectFromAuth(c *gin.Context) (int64, error) {
    auth := c.GetHeader("Authorization")
    if len(auth) < 8 || auth[:7] != "Bearer " {
        return 0, jwt.ErrTokenUnverifiable
    }
    tokenString := auth[7:]
    secret := os.Getenv("JWT_SECRET")
    if secret == "" { secret = "dev_secret_key" }
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil || !token.Valid {
        return 0, jwt.ErrTokenInvalidClaims
    }
    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if !ok { return 0, jwt.ErrTokenInvalidClaims }
    id, err := strconv.ParseInt(claims.Subject, 10, 64)
    if err != nil { return 0, err }
    return id, nil
}
