package atom_news

import (
	atom "BACKEND_SEJUTA_BERITA/atom/news"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PostView: POST /news/:id/view - increments view counter (no auth required)
func PostView(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    if id <= 0 { c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id"}); return }
    if err := atom.RecordViewUseCase(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
        return
    }
    c.Status(http.StatusNoContent)
}
