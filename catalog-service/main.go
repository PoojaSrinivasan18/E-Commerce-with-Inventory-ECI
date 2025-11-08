package main

import (
	catalog_service "github.com/PoojaSrinivasan18/catalog-service/catalog-service"
	"github.com/PoojaSrinivasan18/catalog-service/common"
	"github.com/PoojaSrinivasan18/catalog-service/database"
	"github.com/PoojaSrinivasan18/catalog-service/model"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func main() {

	log.Info("Starting Catalog Service")

	err := common.ConfigSetup("config/dbconfig.yaml")
	if err != nil {
		log.Errorf("ConfigSetup failed: %v", err)
		return
	}

	configuration := common.GetConfig()
	log.Info("Configuration loaded successfully")

	err = database.SetupDB(configuration)
	if err != nil {
		log.Errorf("SetupDB failed: %v", err)
		return
	}

	log.Info("DB Setup Success")

	log.Infof(" Running AutoMigrate...")
	database.GetDB().Exec("SET search_path TO product;")
	err = database.GetDB().AutoMigrate(&model.ProductModel{})
	if err != nil {
		log.Errorf("AutoMigrate failed: %v", err)
	} else {
		log.Infof(" Migration successful!")
	}

	router := gin.Default()
	router.GET("/api/getproductbyid", catalog_service.GetProductById)
	router.GET("/api/getallproduct", catalog_service.GetAllProducts)
	router.POST("/api/addproduct", catalog_service.AddProduct)
	router.DELETE("/api/deleteproduct", catalog_service.DeleteProduct)
	router.PATCH("/api/updateproduct", catalog_service.UpdateProduct)

	router.Run(":3000")

}
