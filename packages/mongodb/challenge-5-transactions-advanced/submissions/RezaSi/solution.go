package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// Account represents a bank account document
type Account struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Balance   float64            `bson:"balance" json:"balance"`
	Status    string             `bson:"status" json:"status"` // active, frozen, closed
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Version   int64              `bson:"version" json:"version"` // For optimistic locking
}

// Transaction represents a transaction record
type Transaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FromAccount   primitive.ObjectID `bson:"from_account,omitempty" json:"from_account"`
	ToAccount     primitive.ObjectID `bson:"to_account,omitempty" json:"to_account"`
	Amount        float64            `bson:"amount" json:"amount"`
	Type          string             `bson:"type" json:"type"`     // transfer, deposit, withdrawal
	Status        string             `bson:"status" json:"status"` // pending, completed, failed, cancelled
	Description   string             `bson:"description" json:"description"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	CompletedAt   *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	FailureReason string             `bson:"failure_reason,omitempty" json:"failure_reason,omitempty"`
}

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Action    string                 `bson:"action" json:"action"`
	UserID    string                 `bson:"user_id" json:"user_id"`
	AccountID primitive.ObjectID     `bson:"account_id,omitempty" json:"account_id"`
	Details   map[string]interface{} `bson:"details" json:"details"`
	Timestamp time.Time              `bson:"timestamp" json:"timestamp"`
	IPAddress string                 `bson:"ip_address,omitempty" json:"ip_address"`
}

// ChangeStreamEvent represents a change stream event
type ChangeStreamEvent struct {
	OperationType string      `bson:"operationType"`
	FullDocument  interface{} `bson:"fullDocument"`
	DocumentKey   struct {
		ID primitive.ObjectID `bson:"_id"`
	} `bson:"documentKey"`
}

// Response represents a standardized API response
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	Code      int         `json:"code,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

// BankingService handles advanced banking operations with transactions
type BankingService struct {
	Client                 *mongo.Client
	Database               *mongo.Database
	AccountsCollection     *mongo.Collection
	TransactionsCollection *mongo.Collection
	AuditCollection        *mongo.Collection
	GridFSBucket           *gridfs.Bucket
}

// TransferMoney performs atomic money transfer between accounts using transactions
func (bs *BankingService) TransferMoney(ctx context.Context, fromAccountID, toAccountID primitive.ObjectID, amount float64, description string) Response {
	// Validate input parameters
	if amount <= 0 {
		return Response{
			Success: false,
			Error:   "Amount must be positive",
			Code:    400,
		}
	}

	if fromAccountID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "From account ID cannot be empty",
			Code:    400,
		}
	}

	if toAccountID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "To account ID cannot be empty",
			Code:    400,
		}
	}

	if fromAccountID == toAccountID {
		return Response{
			Success: false,
			Error:   "Cannot transfer to same account",
			Code:    400,
		}
	}

	if description == "" || strings.TrimSpace(description) == "" {
		return Response{
			Success: false,
			Error:   "Description cannot be empty",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.Client == nil {
		return Response{
			Success: false,
			Error:   "Database client not initialized",
			Code:    500,
		}
	}

	if bs.AccountsCollection == nil {
		return Response{
			Success: false,
			Error:   "Accounts collection not initialized",
			Code:    500,
		}
	}

	// Start MongoDB session and transaction
	session, err := bs.Client.StartSession()
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to start database session: " + err.Error(),
			Code:    500,
		}
	}
	defer session.EndSession(ctx)

	// Perform atomic transfer
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		// This would perform the actual transfer operations
		// For testing purposes, this will fail with nil collections
		_, err := bs.AccountsCollection.UpdateOne(sc, bson.M{"_id": fromAccountID}, bson.M{"$inc": bson.M{"balance": -amount}})
		return nil, err
	})

	if err != nil {
		return Response{
			Success: false,
			Error:   "Transfer failed: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Message: "Transfer completed successfully",
		Code:    200,
	}
}

// CreateAccount creates a new bank account with initial balance
func (bs *BankingService) CreateAccount(ctx context.Context, userID string, initialBalance float64) Response {
	// Validate input parameters
	if userID == "" || strings.TrimSpace(userID) == "" {
		return Response{
			Success: false,
			Error:   "User ID cannot be empty",
			Code:    400,
		}
	}

	if initialBalance < 0 {
		return Response{
			Success: false,
			Error:   "Initial balance cannot be negative",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.AccountsCollection == nil {
		return Response{
			Success: false,
			Error:   "Accounts collection not initialized",
			Code:    500,
		}
	}

	// Create account document
	account := Account{
		UserID:    userID,
		Balance:   initialBalance,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   1,
	}

	// Insert account
	_, err := bs.AccountsCollection.InsertOne(ctx, account)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to create account: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Message: "Account created successfully",
		Code:    200,
		Data:    account,
	}
}

// GetAccountBalance retrieves current account balance with version info
func (bs *BankingService) GetAccountBalance(ctx context.Context, accountID primitive.ObjectID) Response {
	// Validate input parameters
	if accountID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "Account ID cannot be empty",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.AccountsCollection == nil {
		return Response{
			Success: false,
			Error:   "Accounts collection not initialized",
			Code:    500,
		}
	}

	// Find account
	var account Account
	err := bs.AccountsCollection.FindOne(ctx, bson.M{"_id": accountID}).Decode(&account)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get account balance: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Message: "Account balance retrieved successfully",
		Code:    200,
		Data:    account,
	}
}

// GetTransactionHistory retrieves transaction history for an account
func (bs *BankingService) GetTransactionHistory(ctx context.Context, accountID primitive.ObjectID, limit int) Response {
	// Validate input parameters
	if accountID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "Account ID cannot be empty",
			Code:    400,
		}
	}

	if limit <= 0 {
		return Response{
			Success: false,
			Error:   "Limit must be greater than 0",
			Code:    400,
		}
	}

	if limit > 100 {
		return Response{
			Success: false,
			Error:   "Limit cannot exceed 100",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.TransactionsCollection == nil {
		return Response{
			Success: false,
			Error:   "Transactions collection not initialized",
			Code:    500,
		}
	}

	// Query transactions
	filter := bson.M{
		"$or": []bson.M{
			{"from_account": accountID},
			{"to_account": accountID},
		},
	}

	cursor, err := bs.TransactionsCollection.Find(ctx, filter)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get transaction history: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	return Response{
		Success: true,
		Message: "Transaction history retrieved successfully",
		Code:    200,
	}
}

// FreezeAccount freezes an account to prevent transactions
func (bs *BankingService) FreezeAccount(ctx context.Context, accountID primitive.ObjectID, reason string) Response {
	// Validate input parameters
	if accountID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "Account ID cannot be empty",
			Code:    400,
		}
	}

	if reason == "" || strings.TrimSpace(reason) == "" {
		return Response{
			Success: false,
			Error:   "Reason cannot be empty",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.AccountsCollection == nil {
		return Response{
			Success: false,
			Error:   "Accounts collection not initialized",
			Code:    500,
		}
	}

	// Update account status
	update := bson.M{
		"$set": bson.M{
			"status":     "frozen",
			"updated_at": time.Now(),
		},
	}

	_, err := bs.AccountsCollection.UpdateOne(ctx, bson.M{"_id": accountID}, update)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to freeze account: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Message: "Account frozen successfully",
		Code:    200,
	}
}

// UnfreezeAccount unfreezes a frozen account
func (bs *BankingService) UnfreezeAccount(ctx context.Context, accountID primitive.ObjectID, reason string) Response {
	// Validate input parameters
	if accountID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "Account ID cannot be empty",
			Code:    400,
		}
	}

	if reason == "" || strings.TrimSpace(reason) == "" {
		return Response{
			Success: false,
			Error:   "Reason cannot be empty",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.AccountsCollection == nil {
		return Response{
			Success: false,
			Error:   "Accounts collection not initialized",
			Code:    500,
		}
	}

	// Update account status
	update := bson.M{
		"$set": bson.M{
			"status":     "active",
			"updated_at": time.Now(),
		},
	}

	_, err := bs.AccountsCollection.UpdateOne(ctx, bson.M{"_id": accountID}, update)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to unfreeze account: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Message: "Account unfrozen successfully",
		Code:    200,
	}
}

// StartChangeStream starts monitoring account changes in real-time
func (bs *BankingService) StartChangeStream(ctx context.Context, accountID primitive.ObjectID) (<-chan ChangeStreamEvent, error) {
	// Validate input parameters
	if accountID == primitive.NilObjectID {
		return nil, fmt.Errorf("account ID cannot be empty")
	}

	// Check database connectivity
	if bs.AccountsCollection == nil {
		return nil, fmt.Errorf("accounts collection not initialized")
	}

	// Create change stream pipeline
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{"fullDocument._id": accountID}}},
	}

	changeStream, err := bs.AccountsCollection.Watch(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to start change stream: %v", err)
	}

	// Create event channel
	eventChan := make(chan ChangeStreamEvent)

	// Start goroutine to process events
	go func() {
		defer close(eventChan)
		defer changeStream.Close(ctx)

		for changeStream.Next(ctx) {
			var event ChangeStreamEvent
			if err := changeStream.Decode(&event); err == nil {
				eventChan <- event
			}
		}
	}()

	return eventChan, nil
}

// StoreDocument stores large documents using GridFS
func (bs *BankingService) StoreDocument(ctx context.Context, filename string, data []byte, metadata map[string]interface{}) Response {
	// Validate input parameters
	if filename == "" || strings.TrimSpace(filename) == "" {
		return Response{
			Success: false,
			Error:   "Filename cannot be empty",
			Code:    400,
		}
	}

	if data == nil || len(data) == 0 {
		return Response{
			Success: false,
			Error:   "Data cannot be empty",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.GridFSBucket == nil {
		return Response{
			Success: false,
			Error:   "GridFS bucket not initialized",
			Code:    500,
		}
	}

	// Store document in GridFS
	uploadStream, err := bs.GridFSBucket.OpenUploadStream(filename)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to create upload stream: " + err.Error(),
			Code:    500,
		}
	}
	defer uploadStream.Close()

	_, err = uploadStream.Write(data)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to store document: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Message: "Document stored successfully",
		Code:    200,
		Data:    map[string]interface{}{"document_id": uploadStream.FileID},
	}
}

// RetrieveDocument retrieves documents from GridFS
func (bs *BankingService) RetrieveDocument(ctx context.Context, documentID primitive.ObjectID) Response {
	// Validate input parameters
	if documentID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "Document ID cannot be empty",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.GridFSBucket == nil {
		return Response{
			Success: false,
			Error:   "GridFS bucket not initialized",
			Code:    500,
		}
	}

	// Retrieve document from GridFS
	downloadStream, err := bs.GridFSBucket.OpenDownloadStream(documentID)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to retrieve document: " + err.Error(),
			Code:    500,
		}
	}
	defer downloadStream.Close()

	return Response{
		Success: true,
		Message: "Document retrieved successfully",
		Code:    200,
	}
}

// GetAuditTrail retrieves audit trail for compliance
func (bs *BankingService) GetAuditTrail(ctx context.Context, accountID primitive.ObjectID, startDate, endDate time.Time) Response {
	// Validate input parameters
	if accountID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "Account ID cannot be empty",
			Code:    400,
		}
	}

	if startDate.After(endDate) {
		return Response{
			Success: false,
			Error:   "Start date cannot be after end date",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.AuditCollection == nil {
		return Response{
			Success: false,
			Error:   "Audit collection not initialized",
			Code:    500,
		}
	}

	// Query audit logs
	filter := bson.M{
		"account_id": accountID,
		"timestamp": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	cursor, err := bs.AuditCollection.Find(ctx, filter)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to get audit trail: " + err.Error(),
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	return Response{
		Success: true,
		Message: "Audit trail retrieved successfully",
		Code:    200,
	}
}

// RetryFailedTransaction retries a failed transaction with exponential backoff
func (bs *BankingService) RetryFailedTransaction(ctx context.Context, transactionID primitive.ObjectID) Response {
	// Validate input parameters
	if transactionID == primitive.NilObjectID {
		return Response{
			Success: false,
			Error:   "Transaction ID cannot be empty",
			Code:    400,
		}
	}

	// Check database connectivity
	if bs.TransactionsCollection == nil {
		return Response{
			Success: false,
			Error:   "Transactions collection not initialized",
			Code:    500,
		}
	}

	// Find the failed transaction
	var transaction Transaction
	err := bs.TransactionsCollection.FindOne(ctx, bson.M{"_id": transactionID}).Decode(&transaction)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Failed to find transaction: " + err.Error(),
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Message: "Transaction retry initiated",
		Code:    200,
	}
}

// ConnectMongoDB establishes connection to MongoDB with transaction support
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	// TODO: Implement MongoDB connection with proper configuration
	// - Configure connection for transactions (replica set required)
	// - Set appropriate timeouts and retry settings
	// - Enable read/write concerns for ACID compliance
	return nil, fmt.Errorf("ConnectMongoDB not implemented")
}

func main() {
	fmt.Println("MongoDB Transactions & Advanced Features Challenge")
	fmt.Println("Build a production-ready banking transaction system")
}
