package main

import (
	"go-api/model"
	"log"

	"go-api/controller"

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

		r.POST("/order", controller.SaveOrder(db))
		r.GET("/orders", controller.GetOrders(db))
		r.GET("/orders/:order_id", controller.GetOrdersById(db))

		r.GET("/items", controller.GetItems(db))

	log.Fatal(r.Run(":8080"))
}
