package controllers

import (
	"RMS_deploy/handlers"
	Log "RMS_deploy/log"
	"RMS_deploy/models"
	"RMS_deploy/utils"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// @Summary      Create Restaurant
// @Description  Create Restaurant by Logged in Sub-Admin
// @Tags         Sub-Admin
// @Accept       json
// @Produce      json
// @Param        credentials body models.RestaurantRequestBody true "Signup Credentials"
// @Success      200  "Success"
// @Router       /sub-admin/create/restaurant [post]
func CreateRestaurant(c *gin.Context) {
	log := Log.GetLogger(c)
	zap.L().Info("CreateRestaurant is running")

	var body models.RestaurantRequestBody
	err := c.Bind(&body)
	if err != nil {
		//log.Error(err.Error(), zap.String("err", "please send a valid body"))
		//c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		//	"error": "Please send valid body",
		//})
		utils.ResponseWithError(c, http.StatusBadRequest, err, "Please send valid body")
		return
	}
	log.Info("request body", zap.String("name", body.Name),
		zap.Float64("longitude", body.Longitude),
		zap.Float64("latitude", body.Latitude))

	userDetail := utils.GetUserFromContext(c)

	if userDetail.Role != "sub-admin" {
		//log.Info("current role user not authorized to create restaurant")
		//c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "current role user not authorized to create restaurant"})
		utils.ResponseWithError(c, http.StatusUnauthorized, nil, "current role user not authorized to create restaurant")
		return
	}

	restaurant, err := handlers.CreateRestaurant(c, body.Name, body.Longitude, body.Latitude, userDetail.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusNotFound, err, "restaurant not created")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return
	}

	//c.JSON(200, gin.H{
	//	"user": "Restaurant created successfully",
	//})
	utils.ResponseWithSuccess(c, 200, "success", restaurant)
}

// @Summary      Create Dishes
// @Description  Create Dishes for Restaurants by Logged in Sub-Admin
// @Tags         Sub-Admin
// @Accept       json
// @Produce      json
// @Param        credentials body models.DishRequestBody true "Signup Credentials"
// @Success      200  "Success"
// @Router       /sub-admin/create/dish [post]
func CreateDish(c *gin.Context) {
	log := Log.GetLogger(c)
	zap.L().Info("CreateDish is running")

	var body models.DishRequestBody
	err := c.Bind(&body)
	if err != nil {
		utils.ResponseWithError(c, http.StatusBadRequest, err, "Please send valid body")
		return
	}
	log.Info("request body", zap.String("name", body.Name),
		zap.String("price", body.Price.String()),
		zap.Any("tags", body.Tags))
	zap.Uint("restaurantId", body.RestaurantID)

	userDetail := utils.GetUserFromContext(c)

	if userDetail.Role != "sub-admin" {
		utils.ResponseWithError(c, http.StatusUnauthorized, nil, "current role user not authorized to create dishes")
		return
	}

	dish, err := handlers.CreateDish(c, body.Name, body.Price, body.RestaurantID, userDetail.ID, body.Tags)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusNotFound, err, "dish not created")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return
	}

	utils.ResponseWithSuccess(c, 200, "success", dish)

}

// @Summary      Get All Restaurants for all roles
// @Description  Retrieve all Restaurants for the logged-in user
// @Tags         Basic
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Restaurant "List of Restaurants"
// @Router       /restaurants [get]
func GetAllRestaurants(c *gin.Context) {

	log := Log.GetLogger(c)
	log.Info("GetAllRestaurants is running")
	//userDetail, _ := c.Get(initializers.UserString)

	restaurants, err := handlers.GetAllRest(c)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusNotFound, err, "restaurants not found")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return
	}

	//c.JSON(200, gin.H{
	//	"restaurants": restaurants,
	//})
	utils.ResponseWithSuccess(c, 200, "success", restaurants)
}

// @Summary      Get All Dishes for restaurant
// @Description  Retrieve all Dishes for Restaurants for the logged-in user
// @Tags         Basic
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Dish "List of Dishes or restaurant"
// @Router       /dishes/{id} [get]
func GetAllDishesOfRestaurant(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("GetAllDishesOfRestaurant is running")
	id := c.Param("id")
	//userDetail, _ := c.Get(initializers.UserString)

	dishes, err := handlers.GetAllDish(c, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusNotFound, err, "Dishes not found")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return
	}
	utils.ResponseWithSuccess(c, 200, "success", dishes)
}

// @Summary      Get All Restaurants for Admin
// @Description  Retrieve all Restaurants for the logged-in Admin
// @Tags         Admin
// @Tags         Sub-Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Restaurant "List of Restaurants"
// @Router       /admin/get/restaurants [get]
// @Router       /sub-admin/get/restaurants [get]
func GetAdminRestaurants(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("GetAdminRestaurants is running")

	userDetail := utils.GetUserFromContext(c)
	var restaurants []models.Restaurant
	var err error
	if userDetail.Role == models.Admin {
		restaurants, err = handlers.GetAllRestaurantAdmin(c)
	} else if userDetail.Role == models.SubAdmin {
		restaurants, err = handlers.GetAllRestaurantSubAdmin(c, strconv.Itoa(int(userDetail.ID)))
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusNotFound, err, "Restaurants not found")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return

	}
	utils.ResponseWithSuccess(c, 200, "success", restaurants)
}
