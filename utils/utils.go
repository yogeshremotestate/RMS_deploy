package utils

import (
	"RMS_deploy/initializers"
	Log "RMS_deploy/log"
	"RMS_deploy/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func GetUserFromContext(c *gin.Context) (user models.Users) {

	log := Log.GetLogger(c)
	log.Info("GetUserFromContext ")

	// Retrieve user from context
	userDetail, exists := c.Get(initializers.UserString)
	if !exists {
		log.Info("user not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	// Correct type assertion for models.Users
	user, ok := userDetail.(models.Users)
	if !ok {
		log.Info("invalid user type in context")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user type in context"})
		return
	}

	return user

}

// HaversineDistance calculates the great-circle distance between two points
// on the Earth's surface specified by latitude and longitude.
func HaversineDistance(lat1, lon1, lat2, lon2 float64) string {
	// Convert degrees to radians
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	// Apply the Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// EarthRadius is the Earth's radius in kilometers.
	const EarthRadius = 6371.0
	// Return distance in kilometers
	distance := EarthRadius * c
	return fmt.Sprintf("%.2f km", distance)
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func RunMigrations() error {
	// Database connection string
	time.Sleep(2 * time.Second)
	//databaseURL := "postgres://postgres:yogesh1994@database-1.chff2sgkovbj.us-east-1.rds.amazonaws.com:5432/rms?sslmode=require"
	databaseURL := initializers.ENV.MIGRATIONS_URL
	// Create migration instance
	m, err := migrate.New(
		"file://./migrations", // Path to migrations folder
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("error creating migration instance: %v", err)
	}

	// Apply migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %v", err)
	}

	log.Println("Migrations applied successfully.")
	return nil
}

func ResponseWithSuccess(c *gin.Context, code int, message string, data interface{}) {
	log := Log.GetLogger(c)
	log.Info("Sending success response")

	type successResponse struct {
		Message string      `json:"message"`
		Code    int         `json:"code"`
		Data    interface{} `json:"data"`
	}

	response := successResponse{
		Message: message,
		Code:    code,
		Data:    data,
	}

	c.JSON(code, response)
	log.Info("Response sent successfully", zap.Any("response", response))
}

// ResponseWithError sends a structured error response with logging.
func ResponseWithError(c *gin.Context, code int, err error, msg string) {
	log := Log.GetLogger(c)
	type errorResponse struct {
		Error string `json:"error"`
		Code  int    `json:"code"`
	}
	response := errorResponse{
		Error: msg,
		Code:  code,
	}
	c.JSON(code, response)

	log.Error("Error response sent", zap.String("message", msg), zap.Error(err), zap.Int("code", code))

	if code >= 500 {
		log.Error("Internal server error", zap.Error(err))
	}
}
