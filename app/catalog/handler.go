package catalog

import (
	"net/http"
	"strings"

	"github.com/eya20/hiring_test/app/api"
	"gorm.io/gorm"
)

// Product represents a product in the API response
type Product struct {
	Code     string  `json:"code"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

// ProductDetails represents a product with its variants in the API response
type ProductDetails struct {
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category string    `json:"category"`
	Variants []Variant `json:"variants"`
}

// Variant represents a product variant in the API response
type Variant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

// Response represents the catalog API response
type Response struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}

// CatalogHandler handles HTTP requests for catalog operations
type CatalogHandler struct {
	service CatalogService
}

// NewCatalogHandler creates a new catalog handler
func NewCatalogHandler(service CatalogService) *CatalogHandler {
	return &CatalogHandler{
		service: service,
	}
}

// GetCatalog handles GET requests to the catalog endpoint
func (h *CatalogHandler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetProducts()
	if err != nil {
		// Handle different types of errors
		if err.Error() == "database connection failed" {
			api.ErrorResponse(w, http.StatusServiceUnavailable, api.BuildErrorMessage("Database service is temporarily unavailable: ", err))
			return
		}

		// Generic database error
		api.ErrorResponse(w, http.StatusInternalServerError, api.BuildErrorMessage("Unable to retrieve products at this time: ", err))
		return
	}

	response := Response{
		Products: products,
		Total:    len(products),
	}

	api.OKResponse(w, response)
}

// GetProductDetails handles GET requests to the product details endpoint
func (h *CatalogHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	// Extract product code from URL path
	// For simple path extraction: /catalog/PROD001 -> PROD001
	path := r.URL.Path
	code := strings.TrimPrefix(path, "/catalog/")

	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Product code is required")
		return
	}

	product, err := h.service.GetProductByCode(code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.ErrorResponse(w, http.StatusNotFound, "Product not found")
			return
		}

		// Handle different types of errors
		if err.Error() == "database connection failed" {
			api.ErrorResponse(w, http.StatusServiceUnavailable, api.BuildErrorMessage("Database service is temporarily unavailable: ", err))
			return
		}

		// Generic database error
		api.ErrorResponse(w, http.StatusInternalServerError, api.BuildErrorMessage("Unable to retrieve product: ", err))
		return
	}

	api.OKResponse(w, product)
}
