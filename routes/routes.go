package routes

import (
	atom_auth "BACKEND_SEJUTA_BERITA/atom/auth/controller"

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
	return router
}