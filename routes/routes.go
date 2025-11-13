package routes

import (
	atom_auth "BACKEND_SEJUTA_BERITA/atom/auth/controller"
	atom_news "BACKEND_SEJUTA_BERITA/atom/news/controller"
	atom_user "BACKEND_SEJUTA_BERITA/atom/user/controller"
	"BACKEND_SEJUTA_BERITA/config/database"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// authMiddleware extracts user id & role from Authorization: Bearer <jwt>
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if len(auth) < 8 || auth[:7] != "Bearer " {
			c.Next(); return
		}
		tokenString := auth[7:]
		secret := os.Getenv("JWT_SECRET"); if secret == "" { secret = "dev_secret_key" }
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
		if err != nil || !token.Valid { c.Next(); return }
		claims, ok := token.Claims.(*jwt.RegisteredClaims); if !ok { c.Next(); return }
		c.Set("userID", claims.Subject)
		// fetch role quickly (minimal query) only if needed by downstream
		c.Next()
	}
}

// roleRequired ensures a requester has one of allowed roles stored in users table.
func roleRequired(roles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles { allowed[strings.ToLower(r)] = struct{}{} }
	return func(c *gin.Context) {
		// short-circuit if userID missing
		uidAny, exists := c.Get("userID"); if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"}); c.Abort(); return
		}
		uidStr, _ := uidAny.(string)
		// query user role
		// lightweight inline query (could refactor to user resource_db)
	db := database.PgOpenConnection()
	defer db.Close()
		var role string
		if err := db.QueryRow("SELECT role FROM users WHERE id=$1", uidStr).Scan(&role); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"}); c.Abort(); return
		}
		if _, ok := allowed[strings.ToLower(role)]; !ok {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"}); c.Abort(); return
		}
		c.Set("role", role)
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "PUT", "PATCH", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
	}))
	router.Use(authMiddleware())
	// Serve static uploads (e.g., /uploads/xxx.jpg)
	router.Static("/uploads", "./public/uploads")

	// Auth endpoints
	auth := router.Group("/auth")
	{
		auth.POST("/login", atom_auth.PostLoginUser)
		auth.POST("/register", atom_auth.PostRegisterUser)
		auth.POST("/forgot-password", atom_auth.PostForgotPassword)
		auth.POST("/reset-password", atom_auth.PostResetPassword)
	}
	// User profile endpoints (require Bearer token)
	user := router.Group("/user")
	{
		user.GET("/profile", atom_user.GetProfile)
		user.PUT("/profile", atom_user.PutProfile)
	}
	// Public News endpoints
	news := router.Group("/news")
	{
		news.GET("/top", atom_news.GetTop)
		news.GET("/trending", atom_news.GetTrending)
		news.GET("/search", atom_news.GetSearch)
		news.GET("/:id", atom_news.GetByID)
		news.GET("/:id/comments", atom_news.GetComments)
		news.POST("/:id/comments", atom_news.PostComment)
		news.POST("/:id/view", atom_news.PostView)
		news.GET("", atom_news.GetList)
	}

	// Protected content management (admin/editor) using roleRequired middleware
	manage := router.Group("/manage/news", roleRequired("admin", "editor"))
	{
		manage.POST("", atom_news.PostCreate)              // create draft
		manage.PUT("/:id", atom_news.PutUpdate)            // update draft/published
		manage.POST("/:id/publish", atom_news.PostPublish) // publish
		manage.POST("/:id/archive", atom_news.PostArchive) // archive
		manage.GET("/drafts", atom_news.GetDrafts)         // list drafts
		manage.GET("/mine", atom_news.GetMine)             // list own
	}

	upload := router.Group("/manage", roleRequired("admin", "editor"))
	{
		upload.POST("/upload", atom_news.PostUpload)
	}

	// User management (admins only)
	manageUsers := router.Group("/manage/users", roleRequired("admin"))
	{
		manageUsers.GET("", atom_user.GetUsers)
		manageUsers.PUT("/:id/access", atom_user.PutUserAccess)
	}

	// Categories endpoint
	router.GET("/categories", atom_news.GetCategories)
	return router
}