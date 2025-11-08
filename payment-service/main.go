package main

import (
	"github.com/PoojaSrinivasan18/payment-service/common"
	"github.com/PoojaSrinivasan18/payment-service/database"
	"github.com/PoojaSrinivasan18/payment-service/model"
	payment_service "github.com/PoojaSrinivasan18/payment-service/payment-service"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Info("Starting Payment Service")

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

	//	router := gin.Default()
	// routes...
	//router.Run(":3000")

	log.Infof(" Running AutoMigrate...")
	database.GetDB().Exec("SET search_path TO payment;")
	err = database.GetDB().AutoMigrate(&model.PaymentModel{})
	if err != nil {
		log.Errorf("AutoMigrate failed: %v", err)
	} else {
		log.Infof(" Migration successful!")
	}

	router := gin.Default()

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy", "service": "payment"})
	})

	// API versioning with /v1
	v1 := router.Group("/v1")
	{
		v1.GET("/payments/:id", payment_service.GetPaymentById)
		v1.POST("/payments/charge", payment_service.ChargePayment)
		v1.POST("/payments/:id/refund", payment_service.RefundPayment)
		v1.DELETE("/payments/:id", payment_service.DeletePayment)
	}

	//:: Note: For local testing use below
	//router.Run("localhost:3000")

	//:: For Docker use below
	router.Run(":8002")
}
