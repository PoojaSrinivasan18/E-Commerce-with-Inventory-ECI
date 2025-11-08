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
	// Try to get ID from URL parameter first, then query parameter
	productIdStr := c.Param("id")
	if productIdStr == "" {
		productIdStr = c.Query("productId")
	}

	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		log.Errorf("Invalid product ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID", "message": "Product ID must be a valid integer"})
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

func SearchProducts(c *gin.Context) {
	var products []model.ProductModel
	db := database.GetDB()

	// Get query parameters
	name := c.Query("name")
	category := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	isActive := c.Query("is_active")

	// Build the query
	query := db.Model(&model.ProductModel{})

	if name != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+name+"%")
	}
	if category != "" {
		query = query.Where("LOWER(category) LIKE ?", "%"+category+"%")
	}
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}
	if isActive == "true" {
		query = query.Where("is_active = ?", true)
	} else if isActive == "false" {
		query = query.Where("is_active = ?", false)
	}

	// Execute query with pagination
	limit := 50 // Default limit
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := 0
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		log.Errorf("DB search error %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Database search failed"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
		"limit":    limit,
		"offset":   offset,
	})
}
