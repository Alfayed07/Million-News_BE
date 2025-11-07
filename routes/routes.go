package routes

import (
	atom_auth "BACKEND_SEJUTA_BERITA/atom/auth/controller"
	atom_news "BACKEND_SEJUTA_BERITA/atom/news/controller"
	atom_user "BACKEND_SEJUTA_BERITA/atom/user/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "PUT", "PATCH", "DELETE", "GET", "OPTIONS", "TRACE", "CONNECT"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Origin", "Content-Type", "Content-Length", "Date", "origin", "Origins", "x-requested-with", "access-control-allow-methods", "apikey", "Authorization", "Access-Control-Allow-Credentials", "Accept"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
	}))
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
	// News endpoints
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
	// Categories endpoint
	router.GET("/categories", atom_news.GetCategories)
	return router
}