package atom_user

import (
	atom "BACKEND_SEJUTA_BERITA/atom/user"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetUsers handles GET /manage/users listing for admins.
func GetUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    search := strings.TrimSpace(c.Query("search"))

    res, err := atom.ListUsersUseCase(search, page, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, res)
}

// PutUserAccess handles PUT /manage/users/:id/access for admins.
func PutUserAccess(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil || id <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id"})
        return
    }
    var req atom.UpdateUserAccessRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
        return
    }

    actorAny, _ := c.Get("userID")
    actorStr, _ := actorAny.(string)
    actorID, err := strconv.ParseInt(actorStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
        return
    }

    updated, err := atom.UpdateUserAccessUseCase(id, actorID, req)
    if err != nil {
        switch err {
        case atom.ErrInvalidRole:
            c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        case atom.ErrNoUpdateFields:
            c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        case atom.ErrCannotDeactivateSelf, atom.ErrCannotDowngradeSelf:
            c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
        case atom.ErrUserNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"item": updated})
}
