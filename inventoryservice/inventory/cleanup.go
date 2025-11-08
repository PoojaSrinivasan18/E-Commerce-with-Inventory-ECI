package inventory

import (
	database "inventoryservice/database"
	models "inventoryservice/models"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// CleanupExpiredReservations is a background job that releases expired reservations
func CleanupExpiredReservations() {
	log.Info("Starting reservation cleanup job")

	for {
		db := database.GetDB()

		// Find expired reservations
		var expiredReservations []models.ReservationRecord
		if err := db.Where("status = ? AND expires_at < ?", "RESERVED", time.Now()).Find(&expiredReservations).Error; err != nil {
			log.Errorf("Error finding expired reservations: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		if len(expiredReservations) > 0 {
			log.Infof("Found %d expired reservations to clean up", len(expiredReservations))

			tx := db.Begin()

			for _, reservation := range expiredReservations {
				// Find inventory record
				var inventory models.InventoryModel
				if err := tx.Where("product_id = ? AND warehouse = ?",
					reservation.ProductId, reservation.Warehouse).First(&inventory).Error; err != nil {
					log.Errorf("Inventory record not found for reservation %d: %v", reservation.ID, err)
					continue
				}

				// Release reserved quantity back to available stock
				inventory.Reserved -= reservation.Quantity
				inventory.UpdatedAt = time.Now()

				if err := tx.Save(&inventory).Error; err != nil {
					log.Errorf("Failed to release inventory for reservation %d: %v", reservation.ID, err)
					continue
				}

				// Update reservation status
				reservation.Status = "EXPIRED"
				reservation.UpdatedAt = time.Now()

				if err := tx.Save(&reservation).Error; err != nil {
					log.Errorf("Failed to update reservation %d: %v", reservation.ID, err)
					continue
				}

				log.Infof("Released expired reservation %d: product %d, quantity %d, warehouse %s",
					reservation.ID, reservation.ProductId, reservation.Quantity, reservation.Warehouse)
			}

			tx.Commit()
		}

		// Sleep for 1 minute before next cleanup cycle
		time.Sleep(1 * time.Minute)
	}
}

// StartCleanupJob starts the background cleanup job
func StartCleanupJob() {
	go CleanupExpiredReservations()
	log.Info("Reservation cleanup job started")
}

// GetReservationStatus returns current reservation statistics
func GetReservationStatus(c *gin.Context) {
	db := database.GetDB()

	// Count reservations by status
	var stats struct {
		ActiveReservations   int64 `json:"active_reservations"`
		ExpiredReservations  int64 `json:"expired_reservations"`
		ShippedReservations  int64 `json:"shipped_reservations"`
		ReleasedReservations int64 `json:"released_reservations"`
	}

	db.Model(&models.ReservationRecord{}).Where("status = ?", "RESERVED").Count(&stats.ActiveReservations)
	db.Model(&models.ReservationRecord{}).Where("status = ?", "EXPIRED").Count(&stats.ExpiredReservations)
	db.Model(&models.ReservationRecord{}).Where("status = ?", "SHIPPED").Count(&stats.ShippedReservations)
	db.Model(&models.ReservationRecord{}).Where("status = ?", "RELEASED").Count(&stats.ReleasedReservations)

	// Find reservations expiring in the next hour
	var expiringSoon int64
	db.Model(&models.ReservationRecord{}).Where("status = ? AND expires_at BETWEEN ? AND ?",
		"RESERVED", time.Now(), time.Now().Add(1*time.Hour)).Count(&expiringSoon)

	c.JSON(200, gin.H{
		"reservation_stats":  stats,
		"expiring_in_1_hour": expiringSoon,
		"cleanup_active":     true,
	})
}
