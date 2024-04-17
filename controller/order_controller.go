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

		for _, item := range order.Item {
			if err := db.Create(&item).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		order.ItemIDs = make([]uint, len(order.Item))
		for i, item := range order.Item {
			order.ItemIDs[i] = item.ID
		}

		if err := db.Create(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pedido Salvo", "order_id": order.ID})
	}
}

func GetOrders(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var orders []model.Order
		if err := db.Find(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	}
}