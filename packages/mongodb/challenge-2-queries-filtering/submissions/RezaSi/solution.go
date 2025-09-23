package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Product represents a product document in MongoDB
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Category    string             `bson:"category" json:"category"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	Tags        []string           `bson:"tags" json:"tags"`
	Brand       string             `bson:"brand" json:"brand"`
	Rating      float64            `bson:"rating" json:"rating"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// ProductFilter represents filtering criteria for products
type ProductFilter struct {
	Category   string   `json:"category,omitempty"`
	MinPrice   float64  `json:"min_price,omitempty"`
	MaxPrice   float64  `json:"max_price,omitempty"`
	Brand      string   `json:"brand,omitempty"`
	MinRating  float64  `json:"min_rating,omitempty"`
	InStock    bool     `json:"in_stock,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	SearchTerm string   `json:"search_term,omitempty"`
}

// SortOptions represents sorting criteria
type SortOptions struct {
	Field string `json:"field"` // price, rating, name, created_at
	Order int    `json:"order"` // 1 for ascending, -1 for descending
}

// PaginationOptions represents pagination parameters
type PaginationOptions struct {
	Page  int `json:"page"`  // Page number (1-based)
	Limit int `json:"limit"` // Items per page
}

// PaginatedResponse represents a paginated result
type PaginatedResponse struct {
	Data       []Product `json:"data"`
	Total      int64     `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
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

// GetProductsByCategory retrieves products by category with optional filtering
func (ps *ProductService) GetProductsByCategory(ctx context.Context, category string, filter ProductFilter) Response {
	if category == "" {
		return Response{
			Success: false,
			Error:   "Category cannot be empty",
			Code:    400,
		}
	}

	// Build filter starting with category
	queryFilter := bson.M{"category": category}

	// Add additional filters from ProductFilter struct
	if filter.MinPrice > 0 {
		queryFilter["price"] = bson.M{"$gte": filter.MinPrice}
	}
	if filter.MaxPrice > 0 {
		if existingPrice, ok := queryFilter["price"].(bson.M); ok {
			existingPrice["$lte"] = filter.MaxPrice
		} else {
			queryFilter["price"] = bson.M{"$lte": filter.MaxPrice}
		}
	}

	if filter.Brand != "" {
		queryFilter["brand"] = filter.Brand
	}

	if filter.MinRating > 0 {
		queryFilter["rating"] = bson.M{"$gte": filter.MinRating}
	}

	if filter.InStock {
		queryFilter["stock"] = bson.M{"$gt": 0}
	}

	if len(filter.Tags) > 0 {
		queryFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	// Execute query with Find()
	cursor, err := ps.Collection.Find(ctx, queryFilter)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to query products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode products: " + err.Error(),
			Code:    500,
		}
	}

	if products == nil {
		products = []Product{}
	}

	return Response{
		Success: true,
		Data:    products,
		Message: fmt.Sprintf("Found %d products in category %s", len(products), category),
		Code:    200,
	}
}

// GetProductsByPriceRange retrieves products within a price range
func (ps *ProductService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) Response {
	// Handle edge cases
	if minPrice < 0 {
		return Response{
			Success: false,
			Error:   "Minimum price cannot be negative",
			Code:    400,
		}
	}

	if maxPrice < 0 {
		return Response{
			Success: false,
			Error:   "Maximum price cannot be negative",
			Code:    400,
		}
	}

	if minPrice > maxPrice {
		return Response{
			Success: false,
			Error:   "Minimum price cannot be greater than maximum price",
			Code:    400,
		}
	}

	if minPrice == 0 && maxPrice == 0 {
		return Response{
			Success: false,
			Error:   "Price range cannot be zero",
			Code:    400,
		}
	}

	// Create filter with $gte and $lte operators
	filter := bson.M{
		"price": bson.M{
			"$gte": minPrice,
			"$lte": maxPrice,
		},
	}

	// Execute query and return results sorted by price ascending
	opts := options.Find().SetSort(bson.M{"price": 1})
	cursor, err := ps.Collection.Find(ctx, filter, opts)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to query products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode products: " + err.Error(),
			Code:    500,
		}
	}

	if products == nil {
		products = []Product{}
	}

	return Response{
		Success: true,
		Data:    products,
		Message: fmt.Sprintf("Found %d products in price range $%.2f - $%.2f", len(products), minPrice, maxPrice),
		Code:    200,
	}
}

// SearchProductsByName searches products by name using regex
func (ps *ProductService) SearchProductsByName(ctx context.Context, searchTerm string, caseSensitive bool) Response {
	if searchTerm == "" {
		return Response{
			Success: false,
			Error:   "Search term cannot be empty",
			Code:    400,
		}
	}

	// Check for whitespace-only search terms
	if strings.TrimSpace(searchTerm) == "" {
		return Response{
			Success: false,
			Error:   "Search term cannot be empty",
			Code:    400,
		}
	}

	// Create regex pattern for search
	regexOptions := ""
	if !caseSensitive {
		regexOptions = "i" // case-insensitive
	}

	regex := primitive.Regex{
		Pattern: searchTerm,
		Options: regexOptions,
	}

	// Search in both name and description fields using $or
	filter := bson.M{
		"$or": []bson.M{
			{"name": bson.M{"$regex": regex}},
			{"description": bson.M{"$regex": regex}},
		},
	}

	// Execute query
	cursor, err := ps.Collection.Find(ctx, filter)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to search products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode products: " + err.Error(),
			Code:    500,
		}
	}

	if products == nil {
		products = []Product{}
	}

	return Response{
		Success: true,
		Data:    products,
		Message: fmt.Sprintf("Found %d products matching '%s'", len(products), searchTerm),
		Code:    200,
	}
}

// GetProductsWithPagination retrieves products with pagination support
func (ps *ProductService) GetProductsWithPagination(ctx context.Context, pagination PaginationOptions, filter ProductFilter) Response {
	// Validate pagination parameters
	if pagination.Page < 1 {
		return Response{
			Success: false,
			Error:   "Page number must be greater than 0",
			Code:    400,
		}
	}
	if pagination.Limit < 1 {
		return Response{
			Success: false,
			Error:   "Limit must be greater than 0",
			Code:    400,
		}
	}
	if pagination.Limit > 10000 {
		pagination.Limit = 10000 // Prevent excessive data transfer
	}

	// Build filter from ProductFilter
	queryFilter := bson.M{}

	if filter.Category != "" {
		queryFilter["category"] = filter.Category
	}

	if filter.MinPrice > 0 || filter.MaxPrice > 0 {
		priceFilter := bson.M{}
		if filter.MinPrice > 0 {
			priceFilter["$gte"] = filter.MinPrice
		}
		if filter.MaxPrice > 0 {
			priceFilter["$lte"] = filter.MaxPrice
		}
		queryFilter["price"] = priceFilter
	}

	if filter.Brand != "" {
		queryFilter["brand"] = filter.Brand
	}

	if filter.MinRating > 0 {
		queryFilter["rating"] = bson.M{"$gte": filter.MinRating}
	}

	if filter.InStock {
		queryFilter["stock"] = bson.M{"$gt": 0}
	}

	if len(filter.Tags) > 0 {
		queryFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	if filter.SearchTerm != "" {
		regex := primitive.Regex{Pattern: filter.SearchTerm, Options: "i"}
		queryFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": regex}},
			{"description": bson.M{"$regex": regex}},
		}
	}

	// Count total documents matching filter
	total, err := ps.Collection.CountDocuments(ctx, queryFilter)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to count products: " + err.Error(),
			Code:    500,
		}
	}

	// Calculate skip value and total pages
	skip := (pagination.Page - 1) * pagination.Limit
	totalPages := int(math.Ceil(float64(total) / float64(pagination.Limit)))

	// Execute paginated query
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pagination.Limit)).
		SetSort(bson.M{"created_at": -1}) // Most recent first

	cursor, err := ps.Collection.Find(ctx, queryFilter, opts)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to query products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode products: " + err.Error(),
			Code:    500,
		}
	}

	if products == nil {
		products = []Product{}
	}

	// Return PaginatedResponse with metadata
	paginatedResponse := PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: totalPages,
	}

	return Response{
		Success: true,
		Data:    paginatedResponse,
		Message: fmt.Sprintf("Page %d of %d (total: %d products)", pagination.Page, totalPages, total),
		Code:    200,
	}
}

// GetProductsSorted retrieves products with custom sorting
func (ps *ProductService) GetProductsSorted(ctx context.Context, sortOptions []SortOptions, limit int) Response {
	if len(sortOptions) == 0 {
		return Response{
			Success: false,
			Error:   "At least one sort option is required",
			Code:    400,
		}
	}

	// Build sort document from SortOptions array
	sortDoc := bson.D{}
	validFields := map[string]bool{
		"price":      true,
		"rating":     true,
		"name":       true,
		"created_at": true,
		"stock":      true,
		"brand":      true,
	}

	for _, sortOpt := range sortOptions {
		// Validate sort fields to prevent injection
		if !validFields[sortOpt.Field] {
			return Response{
				Success: false,
				Error:   fmt.Sprintf("Invalid sort field: %s", sortOpt.Field),
				Code:    400,
			}
		}

		// Normalize order value
		order := 1
		if sortOpt.Order < 0 {
			order = -1
		}

		sortDoc = append(sortDoc, bson.E{Key: sortOpt.Field, Value: order})
	}

	// Set up find options
	opts := options.Find().SetSort(sortDoc)
	if limit > 0 {
		if limit > 1000 {
			limit = 1000 // Prevent excessive data transfer
		}
		opts.SetLimit(int64(limit))
	}

	// Execute query with sorting options
	cursor, err := ps.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to query products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode products: " + err.Error(),
			Code:    500,
		}
	}

	if products == nil {
		products = []Product{}
	}

	return Response{
		Success: true,
		Data:    products,
		Message: fmt.Sprintf("Retrieved %d products with custom sorting", len(products)),
		Code:    200,
	}
}

// FilterProducts applies complex multi-field filtering
func (ps *ProductService) FilterProducts(ctx context.Context, filter ProductFilter, projection []string) Response {
	// Build comprehensive filter from all ProductFilter fields
	conditions := []bson.M{}

	if filter.Category != "" {
		conditions = append(conditions, bson.M{"category": filter.Category})
	}

	if filter.MinPrice > 0 || filter.MaxPrice > 0 {
		priceFilter := bson.M{}
		if filter.MinPrice > 0 {
			priceFilter["$gte"] = filter.MinPrice
		}
		if filter.MaxPrice > 0 {
			priceFilter["$lte"] = filter.MaxPrice
		}
		conditions = append(conditions, bson.M{"price": priceFilter})
	}

	if filter.Brand != "" {
		conditions = append(conditions, bson.M{"brand": filter.Brand})
	}

	if filter.MinRating > 0 {
		conditions = append(conditions, bson.M{"rating": bson.M{"$gte": filter.MinRating}})
	}

	if filter.InStock {
		conditions = append(conditions, bson.M{"stock": bson.M{"$gt": 0}})
	}

	// Handle array fields (tags) with $in operator
	if len(filter.Tags) > 0 {
		conditions = append(conditions, bson.M{"tags": bson.M{"$in": filter.Tags}})
	}

	if filter.SearchTerm != "" {
		regex := primitive.Regex{Pattern: filter.SearchTerm, Options: "i"}
		conditions = append(conditions, bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": regex}},
				{"description": bson.M{"$regex": regex}},
			},
		})
	}

	// Combine multiple conditions with $and operator
	var queryFilter bson.M
	if len(conditions) > 0 {
		queryFilter = bson.M{"$and": conditions}
	} else {
		queryFilter = bson.M{}
	}

	// Apply field projection to limit returned data
	opts := options.Find()
	if len(projection) > 0 {
		projectionDoc := bson.M{}
		for _, field := range projection {
			projectionDoc[field] = 1
		}
		opts.SetProjection(projectionDoc)
	}

	// Execute query
	cursor, err := ps.Collection.Find(ctx, queryFilter, opts)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to filter products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode products: " + err.Error(),
			Code:    500,
		}
	}

	if products == nil {
		products = []Product{}
	}

	return Response{
		Success: true,
		Data:    products,
		Message: fmt.Sprintf("Filtered results: %d products", len(products)),
		Code:    200,
	}
}

// GetProductsByTags retrieves products that have any of the specified tags
func (ps *ProductService) GetProductsByTags(ctx context.Context, tags []string, matchAll bool) Response {
	if len(tags) == 0 {
		return Response{
			Success: false,
			Error:   "At least one tag is required",
			Code:    400,
		}
	}

	// Validate that no tags are empty
	for _, tag := range tags {
		if strings.TrimSpace(tag) == "" {
			return Response{
				Success: false,
				Error:   "Tags cannot be empty",
				Code:    400,
			}
		}
	}

	// Clean up tags (trim whitespace)
	cleanTags := []string{}
	for _, tag := range tags {
		cleanTags = append(cleanTags, strings.TrimSpace(tag))
	}

	var filter bson.M
	if matchAll {
		// Use $all operator for "all tags" matching
		filter = bson.M{"tags": bson.M{"$all": cleanTags}}
	} else {
		// Use $in operator for "any tag" matching
		filter = bson.M{"tags": bson.M{"$in": cleanTags}}
	}

	cursor, err := ps.Collection.Find(ctx, filter)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to query products by tags: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode products: " + err.Error(),
			Code:    500,
		}
	}

	if products == nil {
		products = []Product{}
	}

	matchType := "any"
	if matchAll {
		matchType = "all"
	}

	return Response{
		Success: true,
		Data:    products,
		Message: fmt.Sprintf("Found %d products matching %s of tags: %v", len(products), matchType, cleanTags),
		Code:    200,
	}
}

// GetTopRatedProducts retrieves highest rated products in category
func (ps *ProductService) GetTopRatedProducts(ctx context.Context, category string, limit int) Response {
	if limit <= 0 {
		return Response{
			Success: false,
			Error:   "Limit must be greater than 0",
			Code:    400,
		}
	}
	if limit > 10000 {
		limit = 10000 // Prevent excessive data transfer
	}

	// Filter by category if specified and only include products with ratings > 0
	filter := bson.M{"rating": bson.M{"$gt": 0}}
	if category != "" {
		filter["category"] = category
	}

	// Sort by rating in descending order and apply limit
	opts := options.Find().
		SetSort(bson.M{"rating": -1, "name": 1}). // Secondary sort by name for consistency
		SetLimit(int64(limit))

	cursor, err := ps.Collection.Find(ctx, filter, opts)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to query top rated products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var products []Product
	if err = cursor.All(ctx, &products); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode products: " + err.Error(),
			Code:    500,
		}
	}

	if products == nil {
		products = []Product{}
	}

	categoryMsg := "all categories"
	if category != "" {
		categoryMsg = fmt.Sprintf("category '%s'", category)
	}

	return Response{
		Success: true,
		Data:    products,
		Message: fmt.Sprintf("Top %d rated products in %s", len(products), categoryMsg),
		Code:    200,
	}
}

// CountProductsByCategory counts products in each category
func (ps *ProductService) CountProductsByCategory(ctx context.Context) Response {
	// Get distinct categories first
	categories, err := ps.Collection.Distinct(ctx, "category", bson.M{})
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get categories: " + err.Error(),
			Code:    500,
		}
	}

	// Count products in each category
	categoryCounts := make(map[string]int64)
	for _, cat := range categories {
		if category, ok := cat.(string); ok {
			count, err := ps.Collection.CountDocuments(ctx, bson.M{"category": category})
			if err != nil {
				return Response{
					Success: false,
					Error:   fmt.Sprintf("Failed to count products in category %s: %v", category, err),
					Code:    500,
				}
			}
			categoryCounts[category] = count
		}
	}

	return Response{
		Success: true,
		Data:    categoryCounts,
		Message: fmt.Sprintf("Product counts for %d categories", len(categoryCounts)),
		Code:    200,
	}
}

// ConnectMongoDB establishes connection to MongoDB
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// SeedSampleProducts creates sample product data for testing
func (ps *ProductService) SeedSampleProducts(ctx context.Context) error {
	// Sample products for testing
	products := []Product{
		{
			ID:          primitive.NewObjectID(),
			Name:        "iPhone 15 Pro",
			Description: "Latest Apple smartphone with advanced features",
			Category:    "Electronics",
			Price:       999.99,
			Stock:       50,
			Tags:        []string{"smartphone", "apple", "premium"},
			Brand:       "Apple",
			Rating:      4.8,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Samsung Galaxy S24",
			Description: "High-performance Android smartphone",
			Category:    "Electronics",
			Price:       899.99,
			Stock:       75,
			Tags:        []string{"smartphone", "samsung", "android"},
			Brand:       "Samsung",
			Rating:      4.6,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Nike Air Max 270",
			Description: "Comfortable running shoes with air cushioning",
			Category:    "Footwear",
			Price:       129.99,
			Stock:       100,
			Tags:        []string{"shoes", "running", "nike", "comfort"},
			Brand:       "Nike",
			Rating:      4.4,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Adidas Ultraboost 22",
			Description: "Premium running shoes with boost technology",
			Category:    "Footwear",
			Price:       189.99,
			Stock:       60,
			Tags:        []string{"shoes", "running", "adidas", "premium"},
			Brand:       "Adidas",
			Rating:      4.7,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "MacBook Pro M3",
			Description: "Professional laptop with M3 chip",
			Category:    "Electronics",
			Price:       1999.99,
			Stock:       25,
			Tags:        []string{"laptop", "apple", "professional", "m3"},
			Brand:       "Apple",
			Rating:      4.9,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Convert to interface{} slice for InsertMany
	documents := make([]interface{}, len(products))
	for i, product := range products {
		documents[i] = product
	}

	_, err := ps.Collection.InsertMany(ctx, documents)
	return err
}

// Example usage and testing
func main() {
	mongoURI := "mongodb://localhost:27017"

	client, err := ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("product_catalog").Collection("products")
	productService := &ProductService{Collection: collection}

	ctx := context.Background()

	// Clear existing data and seed sample products
	collection.Drop(ctx)
	if err := productService.SeedSampleProducts(ctx); err != nil {
		log.Printf("Failed to seed sample data: %v", err)
	} else {
		fmt.Println("Sample products seeded successfully!")
	}

	// Example queries

	// Get products by category
	fmt.Println("\n=== Products in Electronics ===")
	resp := productService.GetProductsByCategory(ctx, "Electronics", ProductFilter{})
	fmt.Printf("Response: %+v\n", resp)

	// Search products by name
	fmt.Println("\n=== Search for 'iPhone' ===")
	resp = productService.SearchProductsByName(ctx, "iPhone", false)
	fmt.Printf("Response: %+v\n", resp)

	// Get products by price range
	fmt.Println("\n=== Products under $200 ===")
	resp = productService.GetProductsByPriceRange(ctx, 0, 200)
	fmt.Printf("Response: %+v\n", resp)

	// Get paginated products
	fmt.Println("\n=== Paginated Products (Page 1) ===")
	resp = productService.GetProductsWithPagination(ctx, PaginationOptions{Page: 1, Limit: 3}, ProductFilter{})
	fmt.Printf("Response: %+v\n", resp)

	// Get top rated products
	fmt.Println("\n=== Top Rated Products ===")
	resp = productService.GetTopRatedProducts(ctx, "", 3)
	fmt.Printf("Response: %+v\n", resp)
}
