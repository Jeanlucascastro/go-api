package main

import (
	"go-api/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=talosorder port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	db.AutoMigrate(&model.Order{}, &model.Item{})

	r := gin.Default()

	r.POST("/order", func(c *gin.Context) {
		var order model.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pedido Salvo", "order_id": order.ID})
	})

	r.GET("/orders", func(c *gin.Context) {
		var orders []model.Order
		if err := db.Preload("Item").Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	})

	log.Fatal(r.Run(":8080"))
}
