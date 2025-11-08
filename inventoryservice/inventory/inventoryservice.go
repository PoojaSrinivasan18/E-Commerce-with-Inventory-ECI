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
