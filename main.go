package main

import (
	"github.com/salamalfis/Assigment2golang/database"
	"github.com/salamalfis/Assigment2golang/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	database.StartDB()
	r := gin.Default()

	// Mengatur rute API
	r.GET("/orders", getOrders)
	r.POST("/orders", createOrder)
	r.POST("/items", createItem)
	r.PUT("/orders/:id", updateOrder)
	r.DELETE("/orders/:id", deleteOrder)

	// Menjalankan server
	if err := r.Run(":80"); err != nil {
		panic(err)
	}
}

func getOrders(c *gin.Context) {
	var orders []models.Orders

	// Retrieve all orders with associated items from the database
	if err := database.GetDB().Preload("Items").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func createOrder(c *gin.Context) {
	var newOrder models.Orders

	
	if err := c.BindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	if newOrder.OrderedAt.IsZero() {
		newOrder.OrderedAt = time.Now()
	}

	// Simpan pesanan ke database
	if err := database.GetDB().Create(&newOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, newOrder)
}

func createItem(c *gin.Context) {
	var item models.Items
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateItem(item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.GetDB().Create(&item)
	c.JSON(http.StatusCreated, item)
}

func validateItem(item models.Items) error {
	if item.ItemCode == "" {
		return errors.New("item code is required")
	}
	if item.Description == "" {
		return errors.New("description is required")
	}
	if item.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	return nil
}
func updateOrder(c *gin.Context) {
	
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	
	var existingOrder models.Orders
	if err := database.GetDB().First(&existingOrder, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var updatedOrder models.Orders
	if err := c.BindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
 
	existingOrder.CustomerName = updatedOrder.CustomerName
	existingOrder.OrderedAt = updatedOrder.OrderedAt

	if err := database.GetDB().Save(&existingOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, existingOrder)
}
func deleteOrder(c *gin.Context) {
	
	orderIDStr := c.Param("id")
	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}


	var existingOrder models.Orders
	if err := database.GetDB().Preload("Items").First(&existingOrder, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Delete  items
	if err := database.GetDB().Delete(&existingOrder.Items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated items"})
		return
	}

	// Delete the order from the database
	if err := database.GetDB().Delete(&existingOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order and associated items deleted successfully"})
}
