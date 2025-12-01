package routes

import (
	"github.com/gbadegesintestimony/jwt-authentication/controllers"
	"github.com/gbadegesintestimony/jwt-authentication/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	// Public routes
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/forgot-password", controllers.ForgotPassword)
			auth.POST("/reset-password", controllers.ResetPassword)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/change-password", controllers.ChangePassword)
			protected.GET("/me", controllers.GetProfile)
			protected.PUT("/me", controllers.UpdateProfile)
		}
	}
}
