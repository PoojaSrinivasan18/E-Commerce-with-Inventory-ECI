package payment_service

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/PoojaSrinivasan18/payment-service/database"
	"github.com/PoojaSrinivasan18/payment-service/model"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func GetPaymentById(c *gin.Context) {
	// Try to get ID from URL parameter first, then query parameter
	paymentIdStr := c.Param("id")
	if paymentIdStr == "" {
		paymentIdStr = c.Query("paymentId")
	}

	paymentId, err := strconv.Atoi(paymentIdStr)
	if err != nil {
		log.Errorf("Invalid payment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID", "message": "Payment ID must be a valid integer"})
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

func ChargePayment(c *gin.Context) {
	var req model.ChargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	db := database.GetDB()

	// Check for existing payment with same idempotency key
	var existingPayment model.PaymentModel
	if err := db.Where("idempotency_key = ?", req.IdempotencyKey).First(&existingPayment).Error; err == nil {
		// Return existing payment
		c.JSON(http.StatusOK, gin.H{
			"message":    "Payment already processed",
			"payment":    existingPayment,
			"idempotent": true,
		})
		return
	}

	// Process new payment
	payment := model.PaymentModel{
		OrderId:        req.OrderId,
		Amount:         req.Amount,
		CustomerId:     req.CustomerId,
		Method:         req.Method,
		Status:         "PROCESSING",
		IdempotencyKey: req.IdempotencyKey,
		Reference:      generatePaymentReference(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Default method if not specified
	if payment.Method == "" {
		payment.Method = "CREDIT_CARD"
	}

	// Simulate payment processing (replace with actual payment gateway)
	success := simulatePaymentProcessing(payment.Amount, payment.Method)

	if success {
		payment.Status = "COMPLETED"
	} else {
		payment.Status = "FAILED"
	}

	payment.UpdatedAt = time.Now()

	// Save payment record
	if err := db.Create(&payment).Error; err != nil {
		log.Errorf("Failed to save payment: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment processing failed"})
		return
	}

	if payment.Status == "COMPLETED" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment processed successfully",
			"payment": payment,
		})
	} else {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error":   "Payment failed",
			"payment": payment,
		})
	}
}

func RefundPayment(c *gin.Context) {
	paymentIdStr := c.Param("id")
	paymentId, err := strconv.Atoi(paymentIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var req model.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	db := database.GetDB()

	// Find original payment
	var payment model.PaymentModel
	if err := db.First(&payment, paymentId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if payment.Status != "COMPLETED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot refund non-completed payment"})
		return
	}

	// Calculate refund amount
	refundAmount := req.Amount
	if refundAmount <= 0 || refundAmount > payment.Amount {
		refundAmount = payment.Amount
	}

	// Create refund record
	refund := model.PaymentModel{
		OrderId:        payment.OrderId,
		Amount:         -refundAmount, // Negative amount for refund
		CustomerId:     payment.CustomerId,
		Method:         payment.Method,
		Status:         "REFUNDED",
		Reference:      generateRefundReference(payment.Reference),
		IdempotencyKey: payment.IdempotencyKey + "_refund_" + strconv.FormatInt(time.Now().Unix(), 10),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save refund record
	if err := db.Create(&refund).Error; err != nil {
		log.Errorf("Failed to save refund: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Refund processing failed"})
		return
	}

	// Update original payment status if full refund
	if refundAmount == payment.Amount {
		payment.Status = "REFUNDED"
		payment.UpdatedAt = time.Now()
		db.Save(&payment)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Refund processed successfully",
		"refund":           refund,
		"original_payment": payment,
	})
}

// generatePaymentReference creates a unique payment reference
func generatePaymentReference() string {
	return fmt.Sprintf("PAY_%d_%d", time.Now().Unix(), rand.Intn(10000))
}

// generateRefundReference creates a refund reference based on original payment
func generateRefundReference(originalRef string) string {
	return fmt.Sprintf("REF_%s_%d", originalRef, time.Now().Unix())
}

// simulatePaymentProcessing simulates payment gateway processing
func simulatePaymentProcessing(amount float64, method string) bool {
	// Simulate different scenarios based on amount
	if amount <= 0 {
		return false
	}

	// Simulate 95% success rate
	return rand.Float64() < 0.95
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
