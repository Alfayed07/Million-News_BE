package atom_news

import (
	atom "BACKEND_SEJUTA_BERITA/atom/news"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetComments(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    items, err := atom.ListCommentsUseCase(id, 100)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"items": items})
}

type postCommentReq struct { Content string `json:"content"` }

func PostComment(c *gin.Context) {
    newsID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    // Optional auth: require token to attach user id; if missing, return 401
    uid, err := subjectFromAuth(c)
    if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"}); return }
    var req postCommentReq
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"}); return }
    cm, err := atom.AddCommentUseCase(newsID, &uid, req.Content)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, cm)
}

func subjectFromAuth(c *gin.Context) (int64, error) {
    auth := c.GetHeader("Authorization")
    if len(auth) < 8 || auth[:7] != "Bearer " { return 0, jwt.ErrTokenUnverifiable }
    tokenString := auth[7:]
    secret := os.Getenv("JWT_SECRET"); if secret == "" { secret = "dev_secret_key" }
    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
    if err != nil || !token.Valid { return 0, jwt.ErrTokenInvalidClaims }
    claims, ok := token.Claims.(*jwt.RegisteredClaims); if !ok { return 0, jwt.ErrTokenInvalidClaims }
    id, err := strconv.ParseInt(claims.Subject, 10, 64)
    if err != nil { return 0, err }
    return id, nil
}
