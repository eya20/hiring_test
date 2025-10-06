package catalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eya20/hiring_test/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoriesRepository is a mock implementation of CategoriesRepositoryInterface
type MockCategoriesRepository struct {
	mock.Mock
}

func (m *MockCategoriesRepository) GetAllCategories() ([]models.Category, error) {
	args := m.Called()
	return args.Get(0).([]models.Category), args.Error(1)
}

func (m *MockCategoriesRepository) GetCategoryByCode(code string) (models.Category, error) {
	args := m.Called(code)
	return args.Get(0).(models.Category), args.Error(1)
}

func (m *MockCategoriesRepository) CreateCategory(category *models.Category) error {
	args := m.Called(category)
	return args.Error(0)
}

func TestCategoriesHandler_GetCategories_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	expectedDbCategories := []models.Category{
		{
			ID:   1,
			Code: "CATGORY001",
			Name: "Clothing",
		},
		{
			ID:   2,
			Code: "CATGORY002",
			Name: "Shoes",
		},
		{
			ID:   3,
			Code: "CATGORY003",
			Name: "Accessories",
		},
	}

	expectedCategories := []Category{
		{
			Code: "CATGORY001",
			Name: "Clothing",
		},
		{
			Code: "CATGORY002",
			Name: "Shoes",
		},
		{
			Code: "CATGORY003",
			Name: "Accessories",
		},
	}

	mockRepo.On("GetAllCategories").Return(expectedDbCategories, nil)

	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetCategories(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response []Category
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategories, response)

	mockRepo.AssertExpectations(t)
}

func TestCategoriesHandler_GetCategories_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	expectedError := errors.New("database connection failed")
	mockRepo.On("GetAllCategories").Return([]models.Category(nil), expectedError)

	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetCategories(w, req)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Database service is temporarily unavailable")

	mockRepo.AssertExpectations(t)
}

func TestCategoriesHandler_GetCategories_GenericError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	expectedError := errors.New("some other error")
	mockRepo.On("GetAllCategories").Return([]models.Category(nil), expectedError)

	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	// Act
	handler.GetCategories(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Unable to retrieve categories at this time")

	mockRepo.AssertExpectations(t)
}

func TestCategoriesHandler_CreateCategory_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	requestBody := CreateCategoryRequest{
		Code: "CATGORY004",
		Name: "Electronics",
	}

	expectedCategory := Category{
		Code: "CATGORY004",
		Name: "Electronics",
	}

	mockRepo.On("CreateCategory", mock.MatchedBy(func(cat *models.Category) bool {
		return cat.Code == "CATGORY004" && cat.Name == "Electronics"
	})).Return(nil)

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response Category
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, response)

	mockRepo.AssertExpectations(t)
}

func TestCategoriesHandler_CreateCategory_InvalidJSON(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	req := httptest.NewRequest("POST", "/categories", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid JSON format", response["error"])

	mockRepo.AssertExpectations(t)
}

func TestCategoriesHandler_CreateCategory_MissingCode(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	requestBody := CreateCategoryRequest{
		Name: "Electronics",
		// Code is missing
	}

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Code is required", response["error"])

	mockRepo.AssertExpectations(t)
}

func TestCategoriesHandler_CreateCategory_MissingName(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	requestBody := CreateCategoryRequest{
		Code: "CATGORY004",
		// Name is missing
	}

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Name is required", response["error"])

	mockRepo.AssertExpectations(t)
}

func TestCategoriesHandler_CreateCategory_DuplicateCode(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	requestBody := CreateCategoryRequest{
		Code: "CATGORY001", // This already exists
		Name: "Electronics",
	}

	expectedError := errors.New("UNIQUE constraint failed: categories.code")
	mockRepo.On("CreateCategory", mock.Anything).Return(expectedError)

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Category with this code already exists", response["error"])

	mockRepo.AssertExpectations(t)
}

func TestCategoriesHandler_CreateCategory_DatabaseError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoriesRepository)
	handler := NewCategoriesHandler(mockRepo)

	requestBody := CreateCategoryRequest{
		Code: "CATGORY004",
		Name: "Electronics",
	}

	expectedError := errors.New("database connection failed")
	mockRepo.On("CreateCategory", mock.Anything).Return(expectedError)

	reqBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateCategory(w, req)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Database service is temporarily unavailable")

	mockRepo.AssertExpectations(t)
}
