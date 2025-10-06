package models

import (
	"gorm.io/gorm"
)

// ProductsRepositoryInterface defines the contract for product repository operations
type ProductsRepositoryInterface interface {
	GetAllProducts() ([]Product, error)
	GetProductsPaginated(offset, limit int) ([]Product, error)
	GetProductsCount() (int64, error)
	GetProductsPaginatedWithFilters(offset, limit int, category string, priceLt *float64) ([]Product, error)
	GetProductsCountWithFilters(category string, priceLt *float64) (int64, error)
	GetProductByCode(code string, product *Product) error
}

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Category").Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductsPaginated retrieves products with pagination
func (r *ProductsRepository) GetProductsPaginated(offset, limit int) ([]Product, error) {
	var products []Product
	if err := r.db.Preload("Category").Preload("Variants").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductsCount returns the total number of products
func (r *ProductsRepository) GetProductsCount() (int64, error) {
	var count int64
	if err := r.db.Model(&Product{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetProductsPaginatedWithFilters retrieves products with pagination and filtering
func (r *ProductsRepository) GetProductsPaginatedWithFilters(offset, limit int, category string, priceLt *float64) ([]Product, error) {
	var products []Product
	query := r.db.Preload("Category").Preload("Variants")

	// Apply category filter
	if category != "" {
		query = query.Joins("JOIN categories ON products.category_id = categories.id").
			Where("categories.name = ?", category)
	}

	// Apply price filter
	if priceLt != nil {
		query = query.Where("products.price < ?", *priceLt)
	}

	if err := query.Order("products.code ASC").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetProductsCountWithFilters returns the total number of products with filters
func (r *ProductsRepository) GetProductsCountWithFilters(category string, priceLt *float64) (int64, error) {
	var count int64
	query := r.db.Model(&Product{})

	// Apply category filter
	if category != "" {
		query = query.Joins("JOIN categories ON products.category_id = categories.id").
			Where("categories.name = ?", category)
	}

	// Apply price filter
	if priceLt != nil {
		query = query.Where("products.price < ?", *priceLt)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetProductByCode retrieves a product by its code with category and variants
func (r *ProductsRepository) GetProductByCode(code string, product *Product) error {
	if err := r.db.Preload("Category").Preload("Variants").Where("code = ?", code).First(product).Error; err != nil {
		return err
	}
	return nil
}
