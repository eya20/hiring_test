package catalog

import (
	"errors"
	"testing"

	"github.com/eya20/hiring_test/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductsRepository is a mock implementation of ProductsRepositoryInterface
type MockProductsRepository struct {
	mock.Mock
}

func (m *MockProductsRepository) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductsRepository) GetProductsPaginated(offset, limit int) ([]models.Product, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductsRepository) GetProductsCount() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductsRepository) GetProductsPaginatedWithFilters(offset, limit int, category string, priceLt *float64) ([]models.Product, error) {
	args := m.Called(offset, limit, category, priceLt)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductsRepository) GetProductsCountWithFilters(category string, priceLt *float64) (int64, error) {
	args := m.Called(category, priceLt)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProductsRepository) GetProductByCode(code string, product *models.Product) error {
	args := m.Called(code, product)
	return args.Error(0)
}

func TestCatalogService_GetProducts_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductsRepository)
	service := NewCatalogService(mockRepo)

	dbProducts := []models.Product{
		{
			ID:    1,
			Code:  "PROD001",
			Price: decimal.NewFromFloat(29.99),
			Category: models.Category{
				Name: "Clothing",
			},
		},
		{
			ID:    2,
			Code:  "PROD002",
			Price: decimal.NewFromFloat(49.99),
			Category: models.Category{
				Name: "Shoes",
			},
		},
	}

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

	mockRepo.On("GetAllProducts").Return(dbProducts, nil)

	// Act
	result, err := service.GetProducts()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, result)
	mockRepo.AssertExpectations(t)
}

func TestCatalogService_GetProducts_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductsRepository)
	service := NewCatalogService(mockRepo)

	expectedError := errors.New("database connection failed")
	mockRepo.On("GetAllProducts").Return([]models.Product(nil), expectedError)

	// Act
	result, err := service.GetProducts()

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestCatalogService_GetProducts_EmptyResult(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductsRepository)
	service := NewCatalogService(mockRepo)

	mockRepo.On("GetAllProducts").Return([]models.Product{}, nil)

	// Act
	result, err := service.GetProducts()

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Len(t, result, 0)
	mockRepo.AssertExpectations(t)
}

func TestCatalogService_GetProductByCode_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockProductsRepository)
	service := NewCatalogService(mockRepo)

	dbProduct := models.Product{
		ID:    1,
		Code:  "PROD001",
		Price: decimal.NewFromFloat(29.99),
		Category: models.Category{
			Name: "Clothing",
		},
		Variants: []models.Variant{
			{
				Name:  "Small",
				SKU:   "PROD001-S",
				Price: decimal.NewFromFloat(29.99),
			},
			{
				Name:  "Large",
				SKU:   "PROD001-L",
				Price: decimal.Zero, // No price set, should inherit from product
			},
		},
	}

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
				Price: 29.99, // Should inherit from product price
			},
		},
	}

	mockRepo.On("GetProductByCode", "PROD001", mock.AnythingOfType("*models.Product")).Run(func(args mock.Arguments) {
		product := args.Get(1).(*models.Product)
		*product = dbProduct
	}).Return(nil)

	// Act
	result, err := service.GetProductByCode("PROD001")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, result)
	mockRepo.AssertExpectations(t)
}
