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

	router := gin.Default()
	router.POST("/api/addinventory", inventory.AddInventory)
	router.PATCH("/api/updateinventory", inventory.UpdateInventory)
	router.DELETE("/api/deleteinventory", inventory.DeleteInventory)
	router.GET("/api/getinventorybyid", inventory.GetInventoryById)
	router.GET("/api/getallinventory", inventory.GetAllInventory)
	router.POST("/api/loadseeddata", inventory.SeedInventoryDetail)

	//:: Note: For local testing use below
	//router.Run("localhost:3000")

	//:: For Docker use below
	router.Run(":3000")
}
