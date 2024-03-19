package main

import (
	"errors"
	"main/data"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func createCar(id int, name string, model string, color string) data.Car {
	car := data.Car{ID: id, Name: name, Model: model, Color: color}
	return car

}

var carToAppend = []data.Car{}

func about(c *gin.Context) {
	c.String(200, "´<h1>The world of code</h1>")
}

func start(c *gin.Context) {
	c.HTML(200, "myWebpage.html", gin.H{"Name": "GMC Yukon", "Model": "2025", "Color": "Black"})
}

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
	router.LoadHTMLGlob("templates/**")
	router.GET("/", start)
	router.GET("/api/about", about)
	router.GET("/api/cars", apiCars)
	router.GET("/api/cars/:id", apiCarsById)

	router.POST("/api/carsadd", apiCarsAdd)
	router.PUT("/api/cars/:id/update", apiCarsUpdate)

	router.Run(":8081")
}
