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

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "catalog"})
	})

	// API versioning with /v1
	v1 := router.Group("/v1")
	{
		v1.GET("/products/:id", catalog_service.GetProductById)
		v1.GET("/products", catalog_service.GetAllProducts)
		v1.POST("/products", catalog_service.AddProduct)
		v1.DELETE("/products/:id", catalog_service.DeleteProduct)
		v1.PATCH("/products/:id", catalog_service.UpdateProduct)
		v1.GET("/products/search", catalog_service.SearchProducts)
	}

	router.Run(":3000")

}
