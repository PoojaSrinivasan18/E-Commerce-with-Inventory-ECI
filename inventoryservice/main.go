package main

import (
	common "inventoryservice/common"
	database "inventoryservice/database"
	inventory "inventoryservice/inventory"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := common.ConfigSetup("configuration/dbconfig.yaml")
	if err != nil {
		log.Error("ConfigSetup failed")
		return
	}

	configuration := common.GetConfig()
	err = database.SetupDB(configuration)

	if err != nil {
		log.Error("SetupDB failed")
		return
	} else {
		log.Info("DB Setup Success")
	}

	// Start reservation cleanup job
	inventory.StartCleanupJob()

	router := gin.Default()

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "inventory"})
	})

	// API versioning with /v1
	v1 := router.Group("/v1")
	{
		v1.POST("/inventory", inventory.AddInventory)
		v1.PATCH("/inventory/:id", inventory.UpdateInventory)
		v1.DELETE("/inventory/:id", inventory.DeleteInventory)
		v1.GET("/inventory/:id", inventory.GetInventoryById)
		v1.GET("/inventory", inventory.GetAllInventory)
		v1.POST("/inventory/seed", inventory.SeedInventoryDetail)

		// New reservation endpoints as per problem statement
		v1.POST("/inventory/reserve", inventory.ReserveInventory)
		v1.POST("/inventory/release", inventory.ReleaseInventory)
		v1.POST("/inventory/ship", inventory.ShipInventory)
		v1.GET("/inventory/availability/:productId", inventory.CheckAvailability)
		v1.GET("/inventory/reservations/status", inventory.GetReservationStatus)
	}

	//:: Note: For local testing use below
	//router.Run("localhost:3000")

	//:: For Docker use below
	router.Run(":3000")
}
