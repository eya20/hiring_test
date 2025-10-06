package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCatalogService is a mock implementation of CatalogService
type MockCatalogService struct {
	mock.Mock
}

func (m *MockCatalogService) GetProducts() ([]Product, error) {
	args := m.Called()
	return args.Get(0).([]Product), args.Error(1)
}

func (m *MockCatalogService) GetProductsPaginated(offset, limit int) ([]Product, int64, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockCatalogService) GetProductsPaginatedWithFilters(offset, limit int, category string, priceLt *float64) ([]Product, int64, error) {
	args := m.Called(offset, limit, category, priceLt)
	return args.Get(0).([]Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockCatalogService) GetProductByCode(code string) (ProductDetails, error) {
	args := m.Called(code)
	return args.Get(0).(ProductDetails), args.Error(1)
}

func TestCatalogHandler_GetCatalog_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCatalogService)
	handler := NewCatalogHandler(mockService)

	expectedProducts := []Product{
		{
			Code:     "PROD001",
			Price:    29.99,
			Category: "Clothing",
		},
		{
			Code:     "PROD002",
			Price:    49.99,
			Category: "Shoes",
		},
	}

	expectedResponse := Response{
		Products: expectedProducts,
		Total:    len(expectedProducts),
	}

	mockService.On("GetProducts").Return(expectedProducts, nil)

	req := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response Response
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)

	mockService.AssertExpectations(t)
}

func TestCatalogHandler_GetCatalog_DatabaseError(t *testing.T) {
	// Arrange
	mockService := new(MockCatalogService)
	handler := NewCatalogHandler(mockService)

	expectedError := errors.New("database connection failed")
	mockService.On("GetProducts").Return([]Product(nil), expectedError)

	req := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Database service is temporarily unavailable")

	mockService.AssertExpectations(t)
}

func TestCatalogHandler_GetCatalog_GenericError(t *testing.T) {
	// Arrange
	mockService := new(MockCatalogService)
	handler := NewCatalogHandler(mockService)

	expectedError := errors.New("some other error")
	mockService.On("GetProducts").Return([]Product(nil), expectedError)

	req := httptest.NewRequest("GET", "/catalog", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetCatalog(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Unable to retrieve products at this time")

	mockService.AssertExpectations(t)
}

func TestCatalogHandler_GetProductDetails_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCatalogService)
	handler := NewCatalogHandler(mockService)

	expectedProduct := ProductDetails{
		Code:     "PROD001",
		Price:    29.99,
		Category: "Clothing",
		Variants: []Variant{
			{
				Name:  "Small",
				SKU:   "PROD001-S",
				Price: 29.99,
			},
			{
				Name:  "Large",
				SKU:   "PROD001-L",
				Price: 34.99,
			},
		},
	}

	mockService.On("GetProductByCode", "PROD001").Return(expectedProduct, nil)

	req := httptest.NewRequest("GET", "/catalog/PROD001", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetProductDetails(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response ProductDetails
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, response)

	mockService.AssertExpectations(t)
}

func TestCatalogHandler_GetProductDetails_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockCatalogService)
	handler := NewCatalogHandler(mockService)

	mockService.On("GetProductByCode", "INVALID").Return(ProductDetails{}, errors.New("record not found"))

	req := httptest.NewRequest("GET", "/catalog/INVALID", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetProductDetails(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Unable to retrieve product")

	mockService.AssertExpectations(t)
}
