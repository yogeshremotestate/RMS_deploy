package routes

import (
	"RMS_deploy/controllers"
	"RMS_deploy/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitializeRoutes(r *gin.Engine) {
	// Swagger setup
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Public routes
	r.POST("/login", controllers.LoginUser)
	r.GET("/restaurants", controllers.GetAllRestaurants)
	r.GET("/dishes/:id", controllers.GetAllDishesOfRestaurant)

	// Admin routes
	adminRoutes := r.Group("/admin", middleware.AuthValidate, middleware.VerifyAdmin)
	{
		adminRoutes.POST("/create/user", controllers.CreateUser)
		adminRoutes.GET("/sub-admins", controllers.GetSubAdmins)
		adminRoutes.GET("/get/restaurants", controllers.GetAdminRestaurants)
	}

	// Sub-admin routes
	subAdminRoutes := r.Group("/sub-admin", middleware.AuthValidate, middleware.VerifySubAdmin)
	{
		subAdminRoutes.POST("/create/user", controllers.CreateUser)
		subAdminRoutes.POST("/create/restaurant", controllers.CreateRestaurant)
		subAdminRoutes.POST("/create/dish", controllers.CreateDish)
		subAdminRoutes.GET("/get/restaurants", controllers.GetAdminRestaurants)
		subAdminRoutes.GET("/get/users", controllers.GetAdminsUsers)
	}

	// User routes
	userRoutes := r.Group("/user", middleware.AuthValidate, middleware.VerifyUser)
	{
		//userRoutes.POST("/create/user", controllers.CreateUser)
		userRoutes.POST("/address", controllers.AddAddress)
		userRoutes.GET("/address", controllers.GetAllAddr)
		userRoutes.POST("/get-distance/:id", controllers.GetRestDistance)
	}
}
