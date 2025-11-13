package atom_news

import (
	atom "BACKEND_SEJUTA_BERITA/atom/news"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetTop(c *gin.Context) {
    // Allow optional limit query param; default higher to support FE slider pages
    lim, _ := strconv.Atoi(c.DefaultQuery("limit", "24"))
    if lim <= 0 { lim = 6 }
    if lim > 100 { lim = 100 }
    items, err := atom.TopUseCase(lim)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"items": items})
}

func GetTrending(c *gin.Context) {
    items, err := atom.TrendingUseCase(5)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"items": items})
}

func GetList(c *gin.Context) {
    category := c.Query("category")
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    res, err := atom.ListUseCase(category, page, limit)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, res)
}

func GetSearch(c *gin.Context) {
    q := c.Query("q")
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    res, err := atom.SearchUseCase(q, page, limit)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, res)
}

func GetByID(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    item, err := atom.DetailUseCase(id)
    if err != nil { c.JSON(http.StatusNotFound, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"item": item})
}

func GetCategories(c *gin.Context) {
    cats, err := atom.ListCategoriesUseCase()
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"items": cats})
}

// Management Handlers (require role middleware)
type createReq struct {
    CategoryID *int64  `json:"category_id"`
    Title      string  `json:"title"`
    Content    string  `json:"content"`
    Image      string  `json:"image"`
}

func PostCreate(c *gin.Context) {
    var req createReq
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"}); return }
    uidStr, _ := c.Get("userID")
    authorID, _ := strconv.ParseInt(uidStr.(string), 10, 64)
    item, err := atom.CreateDraftUseCase(authorID, req.CategoryID, strings.TrimSpace(req.Title), strings.TrimSpace(req.Content), strings.TrimSpace(req.Image))
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"item": item})
}

type updateReq struct {
    CategoryID *int64  `json:"category_id"`
    Title      *string `json:"title"`
    Content    *string `json:"content"`
    Image      *string `json:"image"`
}

func PutUpdate(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    var req updateReq
    if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"}); return }
    item, err := atom.UpdateNewsUseCase(id, req.CategoryID, req.Title, req.Content, req.Image)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"item": item})
}

func PostPublish(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    item, err := atom.PublishNewsUseCase(id)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"item": item})
}

func PostArchive(c *gin.Context) {
    id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    item, err := atom.ArchiveNewsUseCase(id)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, gin.H{"item": item})
}

func GetDrafts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    res, err := atom.ListDraftsUseCase(page, limit)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, res)
}

func GetMine(c *gin.Context) {
    uidStr, _ := c.Get("userID")
    authorID, _ := strconv.ParseInt(uidStr.(string), 10, 64)
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    res, err := atom.ListByAuthorUseCase(authorID, page, limit)
    if err != nil { c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()}); return }
    c.JSON(http.StatusOK, res)
}
