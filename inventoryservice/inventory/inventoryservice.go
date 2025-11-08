package inventory

import (
	"encoding/csv"
	database "inventoryservice/database"
	models "inventoryservice/models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"gorm.io/gorm"
)

func AddInventory(c *gin.Context) {
	var inventoryModel models.InventoryModel
	err := c.ShouldBind(&inventoryModel)
	if err != nil {
		log.Errorf("FORM binding error %v", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	tx := database.GetDB().Create(&inventoryModel)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error saving data"})
		return
	}

	c.IndentedJSON(http.StatusOK, inventoryModel)
}

func UpdateInventory(c *gin.Context) {
	var inventoryModel models.InventoryModel
	err := c.ShouldBind(&inventoryModel)
	if err != nil {
		log.Errorf("FORM binding error %v", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	var existingInventoryDetail models.InventoryModel
	database := database.GetDB()

	t := database.Where("inventory_id=?", inventoryModel.InventoryId).First(&existingInventoryDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	existingInventoryDetail.ProductId = inventoryModel.ProductId
	existingInventoryDetail.WareHouse = inventoryModel.WareHouse
	existingInventoryDetail.OnHand = inventoryModel.OnHand
	existingInventoryDetail.Reserved = inventoryModel.Reserved
	existingInventoryDetail.UpdatedAt = time.Now()

	log.Infof(existingInventoryDetail.WareHouse)

	tx := database.Model(&existingInventoryDetail).Updates(existingInventoryDetail)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error saving data"})
		return
	}

	c.IndentedJSON(http.StatusOK, existingInventoryDetail)
}

func DeleteInventory(c *gin.Context) {
	inventoryId, err := strconv.Atoi(c.Query("inventoryId"))
	if err != nil {
		log.Errorf("Invalid inventory ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid inventory ID"})
		return
	}

	var existingInventoryDetail models.InventoryModel
	database := database.GetDB()

	t := database.Where("inventory_id=?", inventoryId).First(&existingInventoryDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	tx := database.Model(&existingInventoryDetail).Delete(existingInventoryDetail)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error saving data"})
		return
	}

	c.IndentedJSON(http.StatusOK, "Inventory deleted successfully")
}

func GetInventoryById(c *gin.Context) {
	inventoryId, err := strconv.Atoi(c.Query("inventoryId"))
	if err != nil {
		log.Errorf("Invalid inventory ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid inventory ID"})
		return
	}

	var existingInventoryDetail models.InventoryModel
	database := database.GetDB()

	t := database.Where("inventory_id=?", inventoryId).First(&existingInventoryDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	c.IndentedJSON(http.StatusOK, existingInventoryDetail)
}

func GetAllInventory(c *gin.Context) {
	var inventoryDetails []models.InventoryModel
	database := database.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 10
	offset := (page - 1) * limit

	// var inventoryModel models.InventoryModel
	// err := c.ShouldBind(&inventoryModel)
	// if err != nil {
	// 	log.Errorf("FORM binding error %v", err.Error())
	// 	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
	// 	return
	// }

	// t := database.Where("product_id LIKE ?", "%"+inventoryModel.ProductId+"%").
	// 	Offset(offset).Limit(limit).Find(&inventoryDetails)

	t := database.Offset(offset).Limit(limit).Find(&inventoryDetails)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	c.IndentedJSON(http.StatusOK, inventoryDetails)
}

func SeedInventoryDetail(c *gin.Context) {
	log.Infof("Started cleaning up existing inventory data")

	db := database.GetDB()
	if del := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.InventoryModel{}); del.Error != nil {
		log.Errorf("DB delete error: %v", del.Error)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error clearing inventory table"})
		return
	}

	log.Infof("Cleared existing inventory data")

	csvPath := filepath.Join("seeddata", "eci_inventory.csv")
	f, err := os.Open(csvPath)
	if err != nil {
		log.Errorf("Cannot open seed file: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Cannot open seed file"})
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	records, err := r.ReadAll()
	if err != nil {
		log.Errorf("CSV read error: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "CSV read error"})
		return
	}
	if len(records) < 2 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "CSV contains no data"})
		return
	}

	header := records[0]
	idx := make(map[string]int)
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	inserted := 0

	for ri := 1; ri < len(records); ri++ {
		row := records[ri]
		var m models.InventoryModel

		if v, ok := idx["inventory_id"]; ok && v < len(row) {
			if s := strings.TrimSpace(row[v]); s != "" {
				if id, e := strconv.Atoi(s); e == nil {
					m.InventoryId = id
				}
			}
		}

		if v, ok := idx["product_id"]; ok && v < len(row) {
			if s := strings.TrimSpace(row[v]); s != "" {
				if id, e := strconv.Atoi(s); e == nil {
					m.ProductId = id
				}
			}

		}

		if v, ok := idx["warehouse"]; ok && v < len(row) {
			m.WareHouse = strings.TrimSpace(row[v])
		}

		if v, ok := idx["on_hand"]; ok && v < len(row) {
			if s := strings.TrimSpace(row[v]); s != "" {
				if n, e := strconv.Atoi(s); e == nil {
					m.OnHand = n
				}
			}
		}

		if v, ok := idx["reserved"]; ok && v < len(row) {
			if s := strings.TrimSpace(row[v]); s != "" {
				if n, e := strconv.Atoi(s); e == nil {
					m.Reserved = n
				}
			}
		}

		if v, ok := idx["updated_at"]; ok && v < len(row) {
			if s := strings.TrimSpace(row[v]); s != "" {
				// try common timestamp layouts
				var parsed time.Time
				var perr error

				parsed, perr = time.Parse(time.RFC3339, s)
				// if perr != nil {
				// 	parsed, perr = time.Parse("2006-01-02 15:04:05", s)
				// }

				// if perr != nil {
				// 	parsed, perr = time.Parse("2006-01-02", s)
				// }

				if perr == nil {
					m.UpdatedAt = parsed
				} else {
					// fallback
					m.UpdatedAt = time.Now()
				}
			} else {
				m.UpdatedAt = time.Now()
			}
		} else {
			m.UpdatedAt = time.Now()
		}

		tx := db.Create(&m)
		if tx.Error != nil {
			log.Errorf("DB insert error at CSV row %d: %v", ri+1, tx.Error)
			continue
		}
		inserted++
	}

	c.IndentedJSON(http.StatusOK, gin.H{"inserted": inserted})
}

// ReserveInventory reserves inventory for an order with TTL (15 minutes)
func ReserveInventory(c *gin.Context) {
	var req models.ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	db := database.GetDB()

	// Check for duplicate reservation with same idempotency key
	var existingReservation models.ReservationRecord
	if err := db.Where("idempotency_key = ?", req.IdempotencyKey).First(&existingReservation).Error; err == nil {
		// Return existing reservation
		c.JSON(http.StatusOK, gin.H{
			"message":     "Reservation already exists",
			"reservation": existingReservation,
			"idempotent":  true,
		})
		return
	}

	// Start transaction for atomic reservation
	tx := db.Begin()

	// Find inventory to reserve from (try specific warehouse first, then any)
	var inventoryItems []models.InventoryModel
	query := "product_id = ? AND (on_hand - reserved) >= ?"
	args := []interface{}{req.ProductId, req.Quantity}

	if req.Warehouse != "" {
		query += " AND ware_house = ?"
		args = append(args, req.Warehouse)
	}
	query += " ORDER BY ware_house, on_hand DESC"

	if err := tx.Where(query, args...).Find(&inventoryItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if len(inventoryItems) == 0 {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{
			"error":      "Insufficient inventory",
			"product_id": req.ProductId,
			"requested":  req.Quantity,
		})
		return
	}

	// Reserve from the first available warehouse with sufficient stock
	var selectedItem *models.InventoryModel
	for i := range inventoryItems {
		if inventoryItems[i].OnHand-inventoryItems[i].Reserved >= req.Quantity {
			selectedItem = &inventoryItems[i]
			break
		}
	}

	if selectedItem == nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{
			"error":      "Insufficient inventory",
			"product_id": req.ProductId,
			"requested":  req.Quantity,
		})
		return
	}

	// Update inventory reserved count
	selectedItem.Reserved += req.Quantity
	selectedItem.UpdatedAt = time.Now()

	if err := tx.Save(selectedItem).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reserve inventory"})
		return
	}

	// Create reservation record with 15-minute TTL
	reservation := models.ReservationRecord{
		ProductId:      req.ProductId,
		Warehouse:      selectedItem.WareHouse,
		Quantity:       req.Quantity,
		OrderId:        req.OrderId,
		IdempotencyKey: req.IdempotencyKey,
		Status:         "RESERVED",
		ReservedAt:     time.Now(),
		ExpiresAt:      time.Now().Add(15 * time.Minute),
		UpdatedAt:      time.Now(),
	}

	if err := tx.Create(&reservation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation record"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":     "Inventory reserved successfully",
		"reservation": reservation,
		"warehouse":   selectedItem.WareHouse,
		"expires_at":  reservation.ExpiresAt,
	})
}

// ReleaseInventory releases reserved inventory back to available stock
func ReleaseInventory(c *gin.Context) {
	var req models.ReleaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	db := database.GetDB()
	tx := db.Begin()

	// Find reservation record
	var reservation models.ReservationRecord
	if err := tx.Where("idempotency_key = ? AND order_id = ? AND status = ?",
		req.IdempotencyKey, req.OrderId, "RESERVED").First(&reservation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found or already processed"})
		return
	}

	// Find inventory record
	var inventory models.InventoryModel
	if err := tx.Where("product_id = ? AND ware_house = ?",
		reservation.ProductId, reservation.Warehouse).First(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Inventory record not found"})
		return
	}

	// Release reserved quantity back to available stock
	inventory.Reserved -= reservation.Quantity
	inventory.UpdatedAt = time.Now()

	if err := tx.Save(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to release inventory"})
		return
	}

	// Update reservation status
	reservation.Status = "RELEASED"
	reservation.UpdatedAt = time.Now()

	if err := tx.Save(&reservation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reservation record"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":           "Inventory released successfully",
		"reservation":       reservation,
		"released_quantity": reservation.Quantity,
	})
}

// ShipInventory marks reserved inventory as shipped
func ShipInventory(c *gin.Context) {
	var req models.ShipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	db := database.GetDB()
	tx := db.Begin()

	// Find reservation record
	var reservation models.ReservationRecord
	if err := tx.Where("idempotency_key = ? AND order_id = ? AND status = ?",
		req.IdempotencyKey, req.OrderId, "RESERVED").First(&reservation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found or already processed"})
		return
	}

	// Find inventory record
	var inventory models.InventoryModel
	if err := tx.Where("product_id = ? AND ware_house = ?",
		reservation.ProductId, reservation.Warehouse).First(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Inventory record not found"})
		return
	}

	// Ship: reduce both on_hand and reserved quantities
	inventory.OnHand -= reservation.Quantity
	inventory.Reserved -= reservation.Quantity
	inventory.UpdatedAt = time.Now()

	if err := tx.Save(&inventory).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ship inventory"})
		return
	}

	// Update reservation status
	reservation.Status = "SHIPPED"
	reservation.UpdatedAt = time.Now()

	if err := tx.Save(&reservation).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reservation record"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":          "Inventory shipped successfully",
		"reservation":      reservation,
		"shipped_quantity": reservation.Quantity,
	})
}

// CheckAvailability checks product availability across warehouses
func CheckAvailability(c *gin.Context) {
	productIdStr := c.Param("productId")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := database.GetDB()

	var inventoryItems []models.InventoryModel
	if err := db.Where("product_id = ?", productId).Find(&inventoryItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	totalAvailable := 0
	totalOnHand := 0
	totalReserved := 0
	warehouses := make([]gin.H, 0)

	for _, item := range inventoryItems {
		available := item.OnHand - item.Reserved
		totalAvailable += available
		totalOnHand += item.OnHand
		totalReserved += item.Reserved

		warehouses = append(warehouses, gin.H{
			"warehouse": item.WareHouse,
			"on_hand":   item.OnHand,
			"reserved":  item.Reserved,
			"available": available,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id":      productId,
		"total_available": totalAvailable,
		"total_on_hand":   totalOnHand,
		"total_reserved":  totalReserved,
		"warehouses":      warehouses,
	})
}
