package controllers

import (
	"RMS_deploy/handlers"
	"RMS_deploy/initializers"
	Log "RMS_deploy/log"
	"RMS_deploy/models"
	"RMS_deploy/utils"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// @Summary      SignUpUser
// @Description  Register User new account
// @Tags         Sub-Admin
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        credentials body models.UserRequestBody true "Signup Credentials"
// @Success      200  "Success"
// @Router       /admin/create/user [post]
// @Router       /sub-admin/create/user [post]
func CreateUser(c *gin.Context) {
	log := Log.GetLogger(c)
	zap.L().Info("SignUpUser is running")

	var body models.UserRequestBody
	errBody := c.Bind(&body)
	if errBody != nil {
		utils.ResponseWithError(c, http.StatusBadRequest, errBody, "Please send valid body")
		return
	}
	log.Info("request body", zap.String("email", body.Email))

	userDetail := utils.GetUserFromContext(c)

	if body.Role == "sub-admin" {
		if userDetail.Role != models.Admin {
			utils.ResponseWithError(c, http.StatusUnauthorized, nil, "Not authorized to create sub-admin")
			return
		}
	} else if body.Role == "user" {
		if userDetail.Role == models.User {
			utils.ResponseWithError(c, http.StatusUnauthorized, nil, "Not authorized to create User")
			return
		}
	} else if body.Role == "admin" {
		utils.ResponseWithError(c, http.StatusUnauthorized, nil, "Not authorized to create Admin")
		return
	}

	user, err := handlers.UserExist(c, body.Email)
	if err != nil {
		// Log and handle the case where a database error occurs
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("No user found with the provided email, proceeding ", zap.String("email", body.Email))
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
			return
		}
	}

	// Check if the user already exists
	if user.ID > 0 {
		//log.Info("User email already exists", zap.String("email", body.Email))
		//c.JSON(http.StatusBadRequest, gin.H{
		//	"error": "User email already exists",
		//})
		//return
		utils.ResponseWithError(c, http.StatusBadRequest, err, "User Email Already Exists")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		//log.Error(err.Error())
		//c.JSON(500, gin.H{
		//	"error": err.Error(),
		//})
		//return
		utils.ResponseWithError(c, http.StatusInternalServerError, err, "Failed to create user")
		return
	}

	Err := handlers.CreateUser(c, body.Email, string(hash), body.Name, body.Role, userDetail.ID)
	if Err != nil {
		//log.Error("Failed to create user", zap.Error(Err))
		//c.JSON(http.StatusInternalServerError, gin.H{
		//	"error":   "Failed to create user",
		//	"details": Err.Error(),
		//})
		utils.ResponseWithError(c, http.StatusInternalServerError, Err, "Failed to create user")
		return
	}

	utils.ResponseWithSuccess(c, 200, "User created successfully", user)
	//c.JSON(200, gin.H{
	//	"user": "User created successfully",
	//})
}

// @Summary      LoginUser
// @Description  Authenticate a user and return a token
// @Tags         Basic
// @Accept       json
// @Produce      json
// @Param        user body models.LoginRequestBody true "Login Credentials"
// @Success      200  "Success"
// @Router       /login [post]
func LoginUser(c *gin.Context) {
	fmt.Println("reached")
	log := Log.GetLogger(c)
	log.Info("LoginUser is running")
	// time.Sleep(time.Second * 5)
	log.Info("LoginUser has running")
	var body models.LoginRequestBody
	err := c.Bind(&body)
	if err != nil {
		utils.ResponseWithError(c, http.StatusBadRequest, err, "Please send valid body")
		return
	}
	log.Info("request body", zap.String("email", body.Email))

	user, err := handlers.UserExist(c, body.Email)

	//if err == sql.ErrNoRows {
	//	log.Error("user not fount")
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//		"error": " user not found",
	//	})
	//	return
	//}
	//if err != nil {
	//	log.Error(err.Error())
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"error": err.Error(),
	//	})
	//	return
	//}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusBadRequest, err, "User not found")
		} else {
			log.Error("Database query failed", zap.Error(err))
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return
	}

	Err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if Err != nil {
		//log.Error(Err.Error())
		//c.JSON(http.StatusBadRequest, gin.H{
		//	"error": "invalid password",
		//})
		utils.ResponseWithError(c, http.StatusBadRequest, err, "Invalid Password")
		return
	}
	// create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(initializers.ENV.SECRET))
	if err != nil {
		//log.Error(err.Error(), zap.String("err", "failed to generate token"))
		//c.JSON(http.StatusBadRequest, gin.H{
		//	"error": "Failed to Generate token ",
		//})
		utils.ResponseWithError(c, http.StatusInternalServerError, err, "Failed to create token")
		return
	}

	//c.JSON(http.StatusOK, gin.H{"token": tokenString})
	utils.ResponseWithSuccess(c, 200, "Login successfully", tokenString)
}

// @Summary      Adding Address
// @Description  Adding Address by Logged in User
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        credentials body models.AddressRequestBody true "Address Body"
// @Success      200  "Success"
// @Router       /user/address [post]
func AddAddress(c *gin.Context) {
	log := Log.GetLogger(c)
	zap.L().Info("CreateDish is running")

	var body models.AddressRequestBody
	err := c.Bind(&body)
	if err != nil {
		utils.ResponseWithError(c, http.StatusBadRequest, err, "Please send valid body")
		return
	}
	log.Info("request body", zap.Float64("longitude", body.Longitude),
		zap.Float64("latitude", body.Latitude))

	userDetail := utils.GetUserFromContext(c)

	if userDetail.Role != models.User {
		utils.ResponseWithError(c, http.StatusUnauthorized, nil, "current role user not authorized to create address")
		return
	}

	err = handlers.AddAddress(c, body.Name, body.Longitude, body.Latitude, userDetail.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusNotFound, err, "Address adding for user failed")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Failed to add address")
		}
		return
	}
	utils.ResponseWithSuccess(c, 200, "Address added successfully", body)

}

// @Summary      Distance to restaurant
// @Description  Retrieve Distance to Restaurant for logged-in user
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  "Success"
// @Router       /user/get-distance/{id} [get]
func GetRestDistance(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("GetRestDistance is running")
	id := c.Param("id")

	var body models.AddressRequestBody
	err := c.Bind(&body)
	if err != nil {

		utils.ResponseWithError(c, http.StatusBadRequest, err, "Please send valid body")
		return
	}
	log.Info("request body", zap.Float64("longitude", body.Longitude),
		zap.Float64("latitude", body.Latitude),
		zap.String("id", id))

	userDetail := utils.GetUserFromContext(c)

	rest, err := handlers.GetRestById(c, strconv.Itoa(int(userDetail.ID)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusBadRequest, err, "restaurant not found")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return
	}
	distance := utils.HaversineDistance(body.Latitude, body.Longitude, rest.Latitude, rest.Longitude)

	utils.ResponseWithSuccess(c, 200, "success", distance)

}

// @Summary      Get All Address for User
// @Description  Retrieve all Address for logged-in user
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Address "List of Address"
// @Router       /user/address [get]
func GetAllAddr(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("GetAllAddr is running")

	userDetail := utils.GetUserFromContext(c)
	addresses, err := handlers.GetAllAddr(c, strconv.Itoa(int(userDetail.ID)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusBadRequest, err, "Addresses not found")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return
	}
	utils.ResponseWithSuccess(c, 200, "success", addresses)

}

// @Summary      Get Sub-Admins for Admin
// @Description  Retrieve all Sub-Admin for logged-in Admin user
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Users "List of Sub-Admins"
// @Router       /admin/sub-admins [get]
func GetSubAdmins(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("GetSubAdmins is running")

	//userDetail := utils.GetUserFromContext(c)
	subAdmins, err := handlers.GetAllSubAdmins(c)
	if err != nil {
		// Handle specific case where no rows are returned
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusBadRequest, err, "No subdomains  not found")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return

	}
	c.JSON(200, gin.H{
		"data": subAdmins,
	})
	utils.ResponseWithSuccess(c, 200, "success", subAdmins)
}

// @Summary      Get users for Admins
// @Description  Retrieve all users for logged-in Admins/Sub-Admins user
// @Tags         Admin
// @Tags         Sub-Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Users "List of Sub-Admins"
// @Router       /sub-admin/get/users [get]
func GetAdminsUsers(c *gin.Context) {
	log := Log.GetLogger(c)
	log.Info("GetAdminsUsers is running")

	userDetail := utils.GetUserFromContext(c)
	var users []models.Users
	var err error
	if userDetail.Role == models.Admin {
		users, err = handlers.GetAllUsersAdmin(c)
	} else if userDetail.Role == models.SubAdmin {
		users, err = handlers.GetAllUsersSubAdmin(c, strconv.Itoa(int(userDetail.ID)))
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.ResponseWithError(c, http.StatusBadRequest, err, "Users not found")
		} else {
			utils.ResponseWithError(c, http.StatusInternalServerError, err, "Internal server error")
		}
		return
	}
	utils.ResponseWithSuccess(c, 200, "success", users)
}
