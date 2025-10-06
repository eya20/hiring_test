package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/eya20/hiring_test/app/api"
	"github.com/eya20/hiring_test/models"
)

// Category represents a category in the API response
type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CategoriesHandler handles HTTP requests for category operations
type CategoriesHandler struct {
	repo models.CategoriesRepositoryInterface
}

// NewCategoriesHandler creates a new categories handler
func NewCategoriesHandler(repo models.CategoriesRepositoryInterface) *CategoriesHandler {
	return &CategoriesHandler{
		repo: repo,
	}
}

// GetCategories handles GET requests to the categories endpoint
func (h *CategoriesHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	dbCategories, err := h.repo.GetAllCategories()
	if err != nil {
		// Handle different types of errors
		if err.Error() == "database connection failed" {
			api.ErrorResponse(w, http.StatusServiceUnavailable, api.BuildErrorMessage("Database service is temporarily unavailable: ", err))
			return
		}

		// Generic database error
		api.ErrorResponse(w, http.StatusInternalServerError, api.BuildErrorMessage("Unable to retrieve categories at this time: ", err))
		return
	}

	// Map database categories to API categories
	categories := make([]Category, len(dbCategories))
	for i, c := range dbCategories {
		categories[i] = Category{
			Code: c.Code,
			Name: c.Name,
		}
	}

	api.OKResponse(w, categories)
}

// CreateCategory handles POST requests to the categories endpoint
func (h *CategoriesHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate required fields
	if req.Code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Code is required")
		return
	}
	if req.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Name is required")
		return
	}

	// Create database category
	dbCategory := models.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.CreateCategory(&dbCategory); err != nil {
		// Handle different types of errors
		if err.Error() == "database connection failed" {
			api.ErrorResponse(w, http.StatusServiceUnavailable, api.BuildErrorMessage("Database service is temporarily unavailable: ", err))
			return
		}

		// Handle unique constraint violation (duplicate code)
		if err.Error() == "UNIQUE constraint failed: categories.code" ||
			err.Error() == "duplicate key value violates unique constraint" {
			api.ErrorResponse(w, http.StatusConflict, "Category with this code already exists")
			return
		}

		// Generic database error
		api.ErrorResponse(w, http.StatusInternalServerError, api.BuildErrorMessage("Unable to create category: ", err))
		return
	}

	// Return created category
	createdCategory := Category{
		Code: dbCategory.Code,
		Name: dbCategory.Name,
	}

	api.OKResponse(w, createdCategory)
}
