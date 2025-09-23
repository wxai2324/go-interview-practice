package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Product represents a product document in MongoDB
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Brand       string             `bson:"brand" json:"brand"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	Rating      float64            `bson:"rating" json:"rating"`
	Tags        []string           `bson:"tags" json:"tags"`
}

// ProductFilter represents filtering options for products
type ProductFilter struct {
	Category  string  `json:"category,omitempty"`
	Brand     string  `json:"brand,omitempty"`
	MinPrice  float64 `json:"min_price,omitempty"`
	MaxPrice  float64 `json:"max_price,omitempty"`
	MinRating float64 `json:"min_rating,omitempty"`
}

// PaginationOptions represents pagination parameters
type PaginationOptions struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// SortOptions represents sorting parameters
type SortOptions struct {
	Field string `json:"field"`
	Order int    `json:"order"` // 1 for ascending, -1 for descending
}

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// ProductService handles product-related database operations
type ProductService struct {
	Collection *mongo.Collection
}

func main() {
	fmt.Println("MongoDB Queries & Filtering Challenge - Wrong Implementation")
}

// GetProductsByCategory retrieves products by category - WRONG IMPLEMENTATION
func (ps *ProductService) GetProductsByCategory(ctx context.Context, category string, filter ProductFilter) Response {
	// WRONG: No validation at all!
	// Should validate category is not empty

	// WRONG: Using wrong method - should use Find
	_, err := ps.Collection.InsertOne(ctx, Product{}) // This makes no sense!
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error", // WRONG: Generic error message
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []Product{}, // WRONG: Always returns empty array
		Message: "Products retrieved successfully",
		Code:    200,
	}
}

// GetProductsByPriceRange retrieves products within a price range - WRONG IMPLEMENTATION
func (ps *ProductService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) Response {
	// WRONG: No validation at all!
	// Should check for negative prices, min > max, zero prices, etc.

	// WRONG: Using wrong collection method
	_, err := ps.Collection.DeleteOne(ctx, map[string]interface{}{"price": minPrice})
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []Product{}, // WRONG: Always returns empty array
		Message: "Products retrieved successfully",
		Code:    200,
	}
}

// SearchProductsByName searches products by name - WRONG IMPLEMENTATION
func (ps *ProductService) SearchProductsByName(ctx context.Context, searchTerm string, caseSensitive bool) Response {
	// WRONG: No validation at all!
	// Should check for empty search terms, whitespace-only terms, etc.

	// WRONG: Ignoring caseSensitive parameter completely
	// WRONG: Using wrong query method
	_, err := ps.Collection.UpdateOne(ctx, map[string]interface{}{"name": searchTerm}, map[string]interface{}{"$set": map[string]interface{}{"searched": true}})
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []Product{}, // WRONG: Always returns empty array
		Message: "Products retrieved successfully",
		Code:    200,
	}
}

// GetProductsWithPagination retrieves products with pagination - WRONG IMPLEMENTATION
func (ps *ProductService) GetProductsWithPagination(ctx context.Context, pagination PaginationOptions, filter ProductFilter) Response {
	// WRONG: No validation at all!
	// Should validate page > 0, limit > 0, etc.

	// WRONG: Ignoring pagination parameters completely
	// WRONG: Using wrong method
	_, err := ps.Collection.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []Product{}, // WRONG: Always returns empty array
		Message: "Products retrieved successfully",
		Code:    200,
	}
}

// GetProductsSorted retrieves products with sorting - WRONG IMPLEMENTATION
func (ps *ProductService) GetProductsSorted(ctx context.Context, sortOptions []SortOptions, limit int) Response {
	// WRONG: No validation at all!
	// Should validate sortOptions not empty, valid field names, valid order values, limit > 0

	return Response{
		Success: false,
		Error:   "GetProductsSorted not implemented", // WRONG: Not implementing required function
		Code:    500,
	}
}

// FilterProducts filters products with projection - WRONG IMPLEMENTATION
func (ps *ProductService) FilterProducts(ctx context.Context, filter ProductFilter, projection []string) Response {
	// WRONG: No validation at all!
	// Should validate filter parameters, handle projection

	return Response{
		Success: false,
		Error:   "FilterProducts not implemented", // WRONG: Not implementing required function
		Code:    500,
	}
}

// GetProductsByTags retrieves products by tags - WRONG IMPLEMENTATION
func (ps *ProductService) GetProductsByTags(ctx context.Context, tags []string, matchAll bool) Response {
	// WRONG: No validation at all!
	// Should validate tags not empty, no empty strings in tags, etc.

	// WRONG: Ignoring matchAll parameter
	// WRONG: Using wrong method
	err := ps.Collection.Drop(ctx) // WRONG: This drops the entire collection!
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []Product{}, // WRONG: Always returns empty array
		Message: "Products retrieved successfully",
		Code:    200,
	}
}

// GetTopRatedProducts retrieves highest rated products - WRONG IMPLEMENTATION
func (ps *ProductService) GetTopRatedProducts(ctx context.Context, category string, limit int) Response {
	// WRONG: No validation at all!
	// Should validate limit > 0

	// WRONG: Using wrong aggregation
	// WRONG: Not sorting by rating at all
	_, err := ps.Collection.Find(ctx, map[string]interface{}{"rating": 1}) // WRONG: This finds products with rating = 1, not highest rated
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []Product{}, // WRONG: Always returns empty array
		Message: "Products retrieved successfully",
		Code:    200,
	}
}

// CountProductsByCategory counts products by category - WRONG IMPLEMENTATION
func (ps *ProductService) CountProductsByCategory(ctx context.Context) Response {
	// WRONG: Using wrong method - should use aggregation
	count, err := ps.Collection.CountDocuments(ctx, map[string]interface{}{})
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    count, // WRONG: Returns total count, not count by category
		Message: "Count retrieved successfully",
		Code:    200,
	}
}

// SeedSampleProducts seeds sample products - WRONG IMPLEMENTATION
func (ps *ProductService) SeedSampleProducts(ctx context.Context) error {
	// WRONG: Not actually seeding any data
	// Should insert sample products
	return fmt.Errorf("SeedSampleProducts not implemented")
}
