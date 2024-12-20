package main

import (
	_ "RMS_deploy/docs"
	"RMS_deploy/initializers"
	Log "RMS_deploy/log"
	"RMS_deploy/routes"
	Utils "RMS_deploy/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// @title RMS APIs
// @version 1.0
// @description Testing Swagger APIs.
// @termsOfService http://google.com/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api/v1
// @schemes http
func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {

	if err := Log.InitializeLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	//defer Log.LogInstance.Sync()
	if err := Utils.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	r := gin.Default()
	r.Use(Log.LoggerMiddleware())
	routes.InitializeRoutes(r)

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//r.Use(Log.LoggerMiddleware())
	//r.POST("/login", controllers.LoginUser)
	//
	//r.GET("/restaurants", controllers.GetAllRestaurants)
	//r.GET("/dishes/:id", controllers.GetAllDishesOfRestaurant)
	//adminRoutes := r.Group("/admin", middleware.AuthValidate, middleware.VerifyAdmin)
	//{
	//	adminRoutes.POST("/create/user", controllers.CreateUser)
	//	adminRoutes.GET("/sub-admins", controllers.GetSubAdmins)
	//	adminRoutes.GET("/get/restaurants", controllers.GetAdminRestaurants)
	//
	//}
	//subAdminRoutes := r.Group("/sub-admin", middleware.AuthValidate, middleware.VerifySubAdmin)
	//{
	//	//subAdminRoutes.POST("/create/user", controllers.CreateUser) // user can not create any user
	//	subAdminRoutes.POST("/create/restaurant", controllers.CreateRestaurant)
	//	subAdminRoutes.POST("/create/dish", controllers.CreateDish)
	//	subAdminRoutes.GET("/get/restaurants", controllers.GetAdminRestaurants)
	//	subAdminRoutes.GET("/get/users", controllers.GetAdminsUsers)
	//}
	//userRoutes := r.Group("/user", middleware.AuthValidate, middleware.VerifyUser)
	//{
	//	userRoutes.POST("/create/user", controllers.CreateUser)
	//	userRoutes.POST("/address", controllers.AddAddress)
	//	userRoutes.GET("/address", controllers.GetAllAddr)
	//	userRoutes.POST("/get-distance/:id", controllers.GetRestDistance)
	//}
	r.Run()

}
