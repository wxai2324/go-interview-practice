package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Order represents an order document
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

// Customer represents a customer document
type Customer struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Email    string             `bson:"email" json:"email"`
	Location string             `bson:"location" json:"location"`
	JoinDate time.Time          `bson:"join_date" json:"join_date"`
}

// SalesByCategory represents sales analytics by category
type SalesByCategory struct {
	Category      string  `bson:"_id" json:"category"`
	TotalSales    float64 `bson:"total_sales" json:"total_sales"`
	OrderCount    int     `bson:"order_count" json:"order_count"`
	AvgOrderValue float64 `bson:"avg_order_value" json:"avg_order_value"`
}

// TopProduct represents top-selling product analytics
type TopProduct struct {
	ProductID     primitive.ObjectID `bson:"_id" json:"product_id"`
	TotalQuantity int                `bson:"total_quantity" json:"total_quantity"`
	TotalRevenue  float64            `bson:"total_revenue" json:"total_revenue"`
	OrderCount    int                `bson:"order_count" json:"order_count"`
}

// MonthlyRevenue represents revenue analytics by month
type MonthlyRevenue struct {
	Year         int     `bson:"year" json:"year"`
	Month        int     `bson:"month" json:"month"`
	TotalRevenue float64 `bson:"total_revenue" json:"total_revenue"`
	OrderCount   int     `bson:"order_count" json:"order_count"`
}

// CustomerAnalytics represents customer behavior analytics
type CustomerAnalytics struct {
	CustomerID    primitive.ObjectID `bson:"_id" json:"customer_id"`
	TotalSpent    float64            `bson:"total_spent" json:"total_spent"`
	OrderCount    int                `bson:"order_count" json:"order_count"`
	AvgOrderValue float64            `bson:"avg_order_value" json:"avg_order_value"`
	FirstOrder    time.Time          `bson:"first_order" json:"first_order"`
	LastOrder     time.Time          `bson:"last_order" json:"last_order"`
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

// GetSalesByCategory calculates total sales per category using aggregation
func (as *AnalyticsService) GetSalesByCategory(ctx context.Context) Response {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"status": "completed"}}},
		{{"$group", bson.M{
			"_id":             "$category",
			"total_sales":     bson.M{"$sum": "$total"},
			"order_count":     bson.M{"$sum": 1},
			"avg_order_value": bson.M{"$avg": "$total"},
		}}},
		{{"$sort", bson.M{"total_sales": -1}}},
	}

	cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to aggregate sales by category: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var results []SalesByCategory
	if err = cursor.All(ctx, &results); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode results: " + err.Error(),
			Code:    500,
		}
	}

	if results == nil {
		results = []SalesByCategory{}
	}

	return Response{
		Success: true,
		Data:    results,
		Message: fmt.Sprintf("Sales analytics for %d categories", len(results)),
		Code:    200,
	}
}

// GetTopSellingProducts finds top-selling products using aggregation
func (as *AnalyticsService) GetTopSellingProducts(ctx context.Context, limit int) Response {
	if limit <= 0 {
		return Response{
			Success: false,
			Error:   "Limit must be greater than 0",
			Code:    400,
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"status": "completed"}}},
		{{"$group", bson.M{
			"_id":            "$product_id",
			"total_quantity": bson.M{"$sum": "$quantity"},
			"total_revenue":  bson.M{"$sum": "$total"},
			"order_count":    bson.M{"$sum": 1},
		}}},
		{{"$sort", bson.M{"total_quantity": -1}}},
		{{"$limit", limit}},
	}

	cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to aggregate top products: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var results []TopProduct
	if err = cursor.All(ctx, &results); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode results: " + err.Error(),
			Code:    500,
		}
	}

	if results == nil {
		results = []TopProduct{}
	}

	return Response{
		Success: true,
		Data:    results,
		Message: fmt.Sprintf("Top %d selling products", len(results)),
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

// GetRevenueByMonth calculates monthly revenue for a given year
func (as *AnalyticsService) GetRevenueByMonth(ctx context.Context, year int) Response {
	if year <= 0 {
		return Response{
			Success: false,
			Error:   "Year must be greater than 0",
			Code:    400,
		}
	}

	if year < 1900 || year >= 2100 {
		return Response{
			Success: false,
			Error:   "Year must be between 1900 and 2099",
			Code:    400,
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"status": "completed",
			"order_date": bson.M{
				"$gte": time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC),
				"$lt":  time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		}}},
		{{"$group", bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$order_date"},
				"month": bson.M{"$month": "$order_date"},
			},
			"total_revenue": bson.M{"$sum": "$total"},
			"order_count":   bson.M{"$sum": 1},
		}}},
		{{"$sort", bson.M{"_id.month": 1}}},
	}

	cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to aggregate revenue by month: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var results []MonthlyRevenue
	if err = cursor.All(ctx, &results); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode monthly revenue results: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    results,
		Message: fmt.Sprintf("Retrieved monthly revenue for %d", year),
		Code:    200,
	}
}

// GetCustomerAnalytics analyzes customer behavior using aggregation
func (as *AnalyticsService) GetCustomerAnalytics(ctx context.Context, limit int) Response {
	if limit <= 0 {
		return Response{
			Success: false,
			Error:   "Limit must be greater than 0",
			Code:    400,
		}
	}

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"status": "completed"}}},
		{{"$group", bson.M{
			"_id":              "$customer_id",
			"total_spent":      bson.M{"$sum": "$total"},
			"order_count":      bson.M{"$sum": 1},
			"avg_order_value":  bson.M{"$avg": "$total"},
			"first_order_date": bson.M{"$min": "$order_date"},
			"last_order_date":  bson.M{"$max": "$order_date"},
		}}},
		{{"$sort", bson.M{"total_spent": -1}}},
		{{"$limit", limit}},
	}

	cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to aggregate customer analytics: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var results []CustomerAnalytics
	if err = cursor.All(ctx, &results); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode customer analytics results: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    results,
		Message: fmt.Sprintf("Retrieved analytics for top %d customers", len(results)),
		Code:    200,
	}
}

// GetProductPerformance analyzes product performance
func (as *AnalyticsService) GetProductPerformance(ctx context.Context) Response {
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"status": "completed"}}},
		{{"$group", bson.M{
			"_id":         "$product_id",
			"total_sales": bson.M{"$sum": "$total"},
			"units_sold":  bson.M{"$sum": "$quantity"},
			"order_count": bson.M{"$sum": 1},
			"avg_price":   bson.M{"$avg": "$price"},
		}}},
		{{"$sort", bson.M{"total_sales": -1}}},
	}

	cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to aggregate product performance: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var results []interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode product performance results: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    results,
		Message: "Retrieved product performance analytics",
		Code:    200,
	}
}

// GetOrderTrends analyzes order trends over time
func (as *AnalyticsService) GetOrderTrends(ctx context.Context, days int) Response {
	if days <= 0 {
		return Response{
			Success: false,
			Error:   "Days must be greater than 0",
			Code:    400,
		}
	}

	startDate := time.Now().AddDate(0, 0, -days)

	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"order_date": bson.M{"$gte": startDate},
		}}},
		{{"$group", bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$order_date"},
				"month": bson.M{"$month": "$order_date"},
				"day":   bson.M{"$dayOfMonth": "$order_date"},
			},
			"total_orders":  bson.M{"$sum": 1},
			"total_revenue": bson.M{"$sum": "$total"},
		}}},
		{{"$sort", bson.M{"_id": 1}}},
	}

	cursor, err := as.OrdersCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to aggregate order trends: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var results []interface{}
	if err = cursor.All(ctx, &results); err != nil {
		return Response{
			Success: false,
			Error:   "Failed to decode order trends results: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    results,
		Message: fmt.Sprintf("Retrieved order trends for last %d days", days),
		Code:    200,
	}
}

// SeedSampleOrders creates sample order data for testing
func (as *AnalyticsService) SeedSampleOrders(ctx context.Context) error {
	orders := []Order{
		{
			ID:         primitive.NewObjectID(),
			CustomerID: primitive.NewObjectID(),
			ProductID:  primitive.NewObjectID(),
			Quantity:   2,
			Price:      999.99,
			Total:      1999.98,
			Category:   "Electronics",
			OrderDate:  time.Now().AddDate(0, -1, 0),
			Status:     "completed",
		},
		{
			ID:         primitive.NewObjectID(),
			CustomerID: primitive.NewObjectID(),
			ProductID:  primitive.NewObjectID(),
			Quantity:   1,
			Price:      129.99,
			Total:      129.99,
			Category:   "Footwear",
			OrderDate:  time.Now().AddDate(0, -1, -5),
			Status:     "completed",
		},
		{
			ID:         primitive.NewObjectID(),
			CustomerID: primitive.NewObjectID(),
			ProductID:  primitive.NewObjectID(),
			Quantity:   3,
			Price:      899.99,
			Total:      2699.97,
			Category:   "Electronics",
			OrderDate:  time.Now().AddDate(0, -2, 0),
			Status:     "completed",
		},
	}

	documents := make([]interface{}, len(orders))
	for i, order := range orders {
		documents[i] = order
	}

	_, err := as.OrdersCollection.InsertMany(ctx, documents)
	return err
}

func main() {
	mongoURI := "mongodb://localhost:27017"

	client, err := ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("ecommerce_analytics").Collection("orders")
	analyticsService := &AnalyticsService{OrdersCollection: collection}

	ctx := context.Background()

	// Seed sample data
	collection.Drop(ctx)
	if err := analyticsService.SeedSampleOrders(ctx); err != nil {
		log.Printf("Failed to seed sample data: %v", err)
	} else {
		fmt.Println("Sample orders seeded successfully!")
	}

	// Example aggregations
	fmt.Println("\n=== Sales by Category ===")
	resp := analyticsService.GetSalesByCategory(ctx)
	fmt.Printf("Response: %+v\n", resp)

	fmt.Println("\n=== Top Selling Products ===")
	resp = analyticsService.GetTopSellingProducts(ctx, 5)
	fmt.Printf("Response: %+v\n", resp)
}
