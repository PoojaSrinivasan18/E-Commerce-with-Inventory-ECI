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
	router.GET("/api/getpaymentbyid", payment_service.GetPaymentById)
	router.POST("/api/makepayment", payment_service.MakePayment)
	router.DELETE("/api/deletepayment", payment_service.DeletePayment)

	//router.POST("/api/loadseeddata", inventory.SeedInventoryDetail)*/

	//:: Note: For local testing use below
	//router.Run("localhost:3000")

	//:: For Docker use below
	router.Run(":3000")
}
