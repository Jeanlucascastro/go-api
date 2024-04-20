package controller

import (
	"go-api/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SaveOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order model.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var items []model.Item
		if err := db.Find(&items, order.ItemIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		print("items ", items[0].Name)
		order.Item = items

		if err := db.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func GetOrders(db *gorm.DB) gin.HandlerFunc {
  return func(c *gin.Context) {
    var orders []model.Order
    // Include related items using Preload or Joins
    if err := db.Preload("Item").Find(&orders).Error; err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    c.JSON(http.StatusOK, orders)
  }
}
