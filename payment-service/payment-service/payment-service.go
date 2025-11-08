package payment_service

import (
	"net/http"
	"strconv"

	"github.com/PoojaSrinivasan18/payment-service/database"
	"github.com/PoojaSrinivasan18/payment-service/model"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func GetPaymentById(c *gin.Context) {
	paymentId, err := strconv.Atoi(c.Query("paymentId"))
	if err != nil {
		log.Errorf("Invalid payment ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid payment ID"})
		return
	}

	var existingPaymentDetail model.PaymentModel
	database := database.GetDB()

	t := database.Where("payment_id=?", paymentId).First(&existingPaymentDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	c.IndentedJSON(http.StatusOK, existingPaymentDetail)
}
func MakePayment(c *gin.Context) {
	var paymentModel model.PaymentModel
	err := c.ShouldBind(&paymentModel)
	if err != nil {
		log.Errorf("FORM binding error %v", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	tx := database.GetDB().Create(&paymentModel)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error making payment"})
		return
	}

	c.IndentedJSON(http.StatusOK, paymentModel)
}
func DeletePayment(c *gin.Context) {
	paymentId, err := strconv.Atoi(c.Query("paymentId"))
	if err != nil {
		log.Errorf("Invalid payment ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid payment ID"})
		return
	}

	var existingPaymentDetail model.PaymentModel
	database := database.GetDB()

	t := database.Where("payment_id=?", paymentId).First(&existingPaymentDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	tx := database.Model(&existingPaymentDetail).Delete(existingPaymentDetail)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error saving payment data"})
		return
	}

	c.IndentedJSON(http.StatusOK, "Payment deleted successfully")
}
