package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Order represents an order document in MongoDB
type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	ProductID  primitive.ObjectID `bson:"product_id" json:"product_id"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	Price      float64            `bson:"price" json:"price"`
	Total      float64            `bson:"total" json:"total"`
	Category   string             `bson:"category" json:"category"`
	OrderDate  time.Time          `bson:"order_date" json:"order_date"`
	Status     string             `bson:"status" json:"status"`
}

// Customer represents a customer document in MongoDB
type Customer struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	JoinDate time.Time          `bson:"join_date" json:"join_date"`
}

// SalesByCategory represents aggregated sales data by category
type SalesByCategory struct {
	ID            string  `bson:"_id" json:"category"`
	TotalSales    float64 `bson:"total_sales" json:"total_sales"`
	OrderCount    int     `bson:"order_count" json:"order_count"`
	AvgOrderValue float64 `bson:"avg_order_value" json:"avg_order_value"`
}

// TopProduct represents top-selling product data
type TopProduct struct {
	ID            primitive.ObjectID `bson:"_id" json:"product_id"`
	TotalQuantity int                `bson:"total_quantity" json:"total_quantity"`
	TotalRevenue  float64            `bson:"total_revenue" json:"total_revenue"`
	OrderCount    int                `bson:"order_count" json:"order_count"`
}

// MonthlyRevenue represents monthly revenue data
type MonthlyRevenue struct {
	ID           interface{} `bson:"_id" json:"id"`
	TotalRevenue float64     `bson:"total_revenue" json:"total_revenue"`
	OrderCount   int         `bson:"order_count" json:"order_count"`
}

// CustomerAnalytics represents customer behavior analytics
type CustomerAnalytics struct {
	ID             primitive.ObjectID `bson:"_id" json:"customer_id"`
	TotalSpent     float64            `bson:"total_spent" json:"total_spent"`
	OrderCount     int                `bson:"order_count" json:"order_count"`
	AvgOrderValue  float64            `bson:"avg_order_value" json:"avg_order_value"`
	FirstOrderDate time.Time          `bson:"first_order_date" json:"first_order_date"`
	LastOrderDate  time.Time          `bson:"last_order_date" json:"last_order_date"`
}

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// AnalyticsService handles analytics operations using aggregation
type AnalyticsService struct {
	OrdersCollection *mongo.Collection
}

func main() {
	fmt.Println("MongoDB Aggregation Pipeline Challenge - Wrong Implementation")
}

// GetSalesByCategory calculates total sales per category - WRONG IMPLEMENTATION
func (as *AnalyticsService) GetSalesByCategory(ctx context.Context) Response {
	// WRONG: No error handling at all!
	// WRONG: Using wrong aggregation - should group by category
	_, err := as.OrdersCollection.Find(ctx, map[string]interface{}{}) // WRONG: Using Find instead of Aggregate
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error", // WRONG: Generic error message
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []SalesByCategory{}, // WRONG: Always returns empty array
		Message: "Sales retrieved successfully",
		Code:    200,
	}
}

// GetTopSellingProducts finds top-selling products - WRONG IMPLEMENTATION
func (as *AnalyticsService) GetTopSellingProducts(ctx context.Context, limit int) Response {
	// WRONG: No validation at all!
	// Should validate limit > 0

	// WRONG: Using wrong aggregation
	// WRONG: Not grouping by product_id at all
	_, err := as.OrdersCollection.CountDocuments(ctx, map[string]interface{}{}) // WRONG: Using CountDocuments instead of Aggregate
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []TopProduct{}, // WRONG: Always returns empty array
		Message: "Top products retrieved successfully",
		Code:    200,
	}
}

// GetRevenueByMonth calculates monthly revenue - WRONG IMPLEMENTATION
func (as *AnalyticsService) GetRevenueByMonth(ctx context.Context, year int) Response {
	// WRONG: No validation at all!
	// Should validate year > 0, year range, etc.

	// WRONG: Using wrong method - should use aggregation with date operators
	_, err := as.OrdersCollection.UpdateMany(ctx, map[string]interface{}{"year": year}, map[string]interface{}{"$set": map[string]interface{}{"processed": true}})
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []MonthlyRevenue{}, // WRONG: Always returns empty array
		Message: "Monthly revenue retrieved successfully",
		Code:    200,
	}
}

// GetCustomerAnalytics analyzes customer behavior - WRONG IMPLEMENTATION
func (as *AnalyticsService) GetCustomerAnalytics(ctx context.Context, limit int) Response {
	// WRONG: No validation at all!
	// Should validate limit > 0

	// WRONG: Using wrong method - should use aggregation to group by customer_id
	_, err := as.OrdersCollection.DeleteMany(ctx, map[string]interface{}{"limit": limit}) // WRONG: This deletes data!
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []CustomerAnalytics{}, // WRONG: Always returns empty array
		Message: "Customer analytics retrieved successfully",
		Code:    200,
	}
}

// GetProductPerformance analyzes product performance - WRONG IMPLEMENTATION
func (as *AnalyticsService) GetProductPerformance(ctx context.Context) Response {
	// WRONG: Using wrong method - should use aggregation
	// WRONG: Not grouping by product at all
	_, err := as.OrdersCollection.InsertOne(ctx, map[string]interface{}{"performance": "calculated"}) // WRONG: This inserts random data!
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []interface{}{}, // WRONG: Always returns empty array
		Message: "Product performance retrieved successfully",
		Code:    200,
	}
}

// GetOrderTrends analyzes order trends - WRONG IMPLEMENTATION
func (as *AnalyticsService) GetOrderTrends(ctx context.Context, days int) Response {
	// WRONG: No validation at all!
	// Should validate days > 0

	// WRONG: Using wrong method - should use aggregation with date filtering
	err := as.OrdersCollection.Drop(ctx) // WRONG: This drops the entire collection!
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    []interface{}{}, // WRONG: Always returns empty array
		Message: "Order trends retrieved successfully",
		Code:    200,
	}
}

// SeedSampleOrders creates sample order data - WRONG IMPLEMENTATION
func (as *AnalyticsService) SeedSampleOrders(ctx context.Context) error {
	// WRONG: Not actually seeding any data
	// Should insert sample orders
	return fmt.Errorf("SeedSampleOrders not implemented")
}
