package main

import (
	"errors"
	"log"
	"main/data"
	"main/rate"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func apiCars(c *gin.Context) {
	var cars []data.Car
	data.DB.Find(&cars)

	c.IndentedJSON(http.StatusOK, cars)
}
func apiCarsById(c *gin.Context) {
	id := c.Param("id")

	var car data.Car
	err := data.DB.First(&car, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		c.IndentedJSON(http.StatusOK, car)
	}
}

func apiCarsAdd(c *gin.Context) {
	var car data.Car

	if err := c.BindJSON(&car); err != nil {
		return
	}

	car.ID = 0
	err := data.DB.Create(&car).Error
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.IndentedJSON(http.StatusCreated, car)
	}

}

func apiCarsUpdate(c *gin.Context) {
	id := c.Param("id")

	var car data.Car
	err := data.DB.First(&car, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		if err := c.BindJSON(&car); err != nil {
			return
		}

		car.ID, _ = strconv.Atoi(id)
		data.DB.Save(&car)
		c.IndentedJSON(http.StatusOK, car)

	}

}

func apiCarsDelete(c *gin.Context) {
	id := c.Param("id")

	var car data.Car
	err := data.DB.First(&car, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	if err != nil {
		// Handle other potential errors
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "database error", "error": err.Error()})
		return
	}

	// Perform the deletion
	result := data.DB.Delete(&car)
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "delete failed", "error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "car deleted"})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the origin of the request
		origin := c.Request.Header.Get("Origin")

		// Set headers to allow the specific origin, replace with your client URL in production
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func login(c *gin.Context) {
	var loginDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login details"})
		return
	}

	var user data.User
	if err := data.DB.Where("username = ?", loginDetails.Username).First(&user).Error; err != nil {
		log.Println("Login error for user:", loginDetails.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"}) // Generic error message
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password)); err != nil {
		log.Println("Password verification failed for user:", loginDetails.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	if err := session.Save(); err != nil {
		log.Println("Error saving session:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func registerUser(c *gin.Context) {
	var registrationDetails struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&registrationDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration details"})
		return
	}

	// Check if username already exists
	var existingUser data.User
	if err := data.DB.Where("username = ?", registrationDetails.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already in use"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registrationDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	// Create user record
	newUser := data.User{
		Username: registrationDetails.Username,
		Password: string(hashedPassword),
	}
	if err := data.DB.Create(&newUser).Error; err != nil {
		log.Println("Failed to create user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func onHealth(c *gin.Context) {
	c.Status(200)
	c.Writer.Write([]byte("ok"))
}

func main() {
	var config Config
	readConfig(&config)
	data.Init(config.Database.File,
		config.Database.Server,
		config.Database.Database,
		config.Database.Username,
		config.Database.Password,
		config.Database.Port)

	//create webpage with route
	router := gin.Default()

	// Initialize the rate limiter
	limiter := rate.NewRateLimiter(1, 5)
	router.Use(CORSMiddleware())

	// Initialize Redis store for session management
	store, err := redis.NewStore(10, "tcp", config.Redis.Server, "", []byte("secret"))
	if err != nil {
		log.Fatal("failed to connect to Redis:", err)
	}
	router.Use(sessions.Sessions("mysession", store))

	router.POST("/register", limiter.Middleware(), registerUser)
	router.POST("/login", limiter.Middleware(), login)
	router.GET("/logout", limiter.Middleware(), logout)

	router.GET("/api/cars", apiCars)
	router.GET("/api/cars/:id", apiCarsById)

	router.POST("/api/carsadd", apiCarsAdd)
	router.PUT("/api/cars/:id/update", apiCarsUpdate)

	router.DELETE("/api/carsdelete/:id", apiCarsDelete)

	router.GET("/healthz", onHealth)

	router.Run(":8080")
}
