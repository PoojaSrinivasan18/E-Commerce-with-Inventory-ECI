package catalog_service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PoojaSrinivasan18/catalog-service/database"
	"github.com/PoojaSrinivasan18/catalog-service/model"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func GetProductById(c *gin.Context) {
	productId, err := strconv.Atoi(c.Query("productId"))
	if err != nil {
		log.Errorf("Invalid product ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var existingProductDetail model.ProductModel
	database := database.GetDB()

	t := database.Where("product_id=?", productId).First(&existingProductDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	c.IndentedJSON(http.StatusOK, existingProductDetail)
}
func GetAllProducts(c *gin.Context) {
	var products []model.ProductModel
	db := database.GetDB()

	t := db.Find(&products)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": t.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}

func AddProduct(c *gin.Context) {
	var productModel model.ProductModel
	err := c.ShouldBind(&productModel)
	if err != nil {
		log.Errorf("FORM binding error %v", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	tx := database.GetDB().Create(&productModel)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error adding product"})
		return
	}

	c.IndentedJSON(http.StatusOK, productModel)
}
func DeleteProduct(c *gin.Context) {
	productId, err := strconv.Atoi(c.Query("productId"))
	if err != nil {
		log.Errorf("Invalid product ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var existingProductDetail model.ProductModel
	database := database.GetDB()

	t := database.Where("product_id=?", productId).First(&existingProductDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	tx := database.Model(&existingProductDetail).Delete(existingProductDetail)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error saving product data"})
		return
	}

	c.IndentedJSON(http.StatusOK, "Product deleted successfully")
}

/*
	func UpdateProduct(c *gin.Context) {
		var product model.ProductModel
		database := database.GetDB()

		// Bind JSON body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
			return
		}

		// Validate product_id
		if product.ProductId == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Product ID is required"})
			return
		}

		var existingProduct model.ProductModel
		// Try to find the product by product_id
		if err := database.First(&existingProduct, "product_id = ?", product.ProductId).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Invalid product ID"})
			return
		}

		// Update fields
		existingProduct.Sku = product.Sku
		existingProduct.Price = product.Price
		existingProduct.Name = product.Name
		existingProduct.Category = product.Category
		existingProduct.IsActive = product.IsActive
		existingProduct.Description = product.Description
		existingProduct.UpdatedAt = time.Now()

		// Save updated product
		if err := database.Save(&existingProduct).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Product updated successfully",
			"product": existingProduct,
		})
	}
*/
func UpdateProduct(c *gin.Context) {
	var product model.ProductModel
	database := database.GetDB()

	// Bind JSON body
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Validate product_id
	if product.ProductId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product ID is required"})
		return
	}

	var existingProduct model.ProductModel
	// Try to find the product by product_id
	if err := database.First(&existingProduct, "product_id = ?", product.ProductId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Invalid product ID"})
		return
	}
	// Update fields
	if product.Sku != "" {
		existingProduct.Sku = product.Sku
	}
	if product.Price != 0.0 {
		existingProduct.Price = product.Price
	}
	if product.Name != "" {
		existingProduct.Name = product.Name
	}

	if product.Category != "" {
		existingProduct.Category = product.Category
	}
	if !product.IsActive {
		existingProduct.IsActive = product.IsActive
	}
	if product.Description != "" {
		existingProduct.Description = product.Description
	}

	existingProduct.UpdatedAt = time.Now()

	// Save updated product
	if err := database.Save(&existingProduct).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": existingProduct,
	})
}
