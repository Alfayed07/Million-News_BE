package atom_news

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// PostUpload handles POST /manage/upload (multipart/form-data with field "file")
// Returns { path: "/uploads/filename.ext" }
func PostUpload(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message":"file not provided"}); return }
    // Basic extension allow list
    ext := strings.ToLower(filepath.Ext(file.Filename))
    allowed := map[string]struct{}{ ".jpg":{}, ".jpeg":{}, ".png":{}, ".webp":{}, ".gif":{} }
    if _, ok := allowed[ext]; !ok { c.JSON(http.StatusBadRequest, gin.H{"message":"unsupported file type"}); return }
    // Sanitize base name
    base := strings.TrimSuffix(file.Filename, ext)
    base = strings.Map(func(r rune) rune {
        if (r>='a'&&r<='z') || (r>='A'&&r<='Z') || (r>='0'&&r<='9') { return r }
        return '-'
    }, base)
    if len(base) == 0 { base = "upload" }
    finalName := base + "-" + time.Now().Format("20060102150405") + ext
    // Ensure uploads dir exists
    upDir := filepath.FromSlash("public/uploads")
    if err := os.MkdirAll(upDir, 0755); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"message": "cannot create uploads dir"}); return }
    dst := filepath.Join(upDir, finalName)
    if err := c.SaveUploadedFile(file, dst); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"message":"failed to save"}); return }
    c.JSON(http.StatusOK, gin.H{"path": "/uploads/" + finalName})
}