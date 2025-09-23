package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User represents a user document in MongoDB
type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Age   int                `bson:"age" json:"age"`
}

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Age   int    `json:"age" bson:"age"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Age   int    `json:"age,omitempty" bson:"age,omitempty"`
}

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// UserService handles user-related database operations
type UserService struct {
	Collection *mongo.Collection
}

func main() {
	fmt.Println("MongoDB CRUD Challenge - Wrong Implementation")
}

// CreateUser creates a new user in the database - WRONG IMPLEMENTATION
func (us *UserService) CreateUser(ctx context.Context, req CreateUserRequest) Response {
	// WRONG: No validation at all!
	// Should validate name, email format, age range, etc.

	user := User{
		ID:    primitive.NewObjectID(),
		Name:  req.Name,  // WRONG: Accepts empty names
		Email: req.Email, // WRONG: No email validation
		Age:   req.Age,   // WRONG: Accepts negative ages, unrealistic ages
	}

	// WRONG: Using wrong method - should use InsertOne
	_, err := us.Collection.Find(ctx, user) // This makes no sense!
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error", // WRONG: Generic error message
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    user,
		Message: "User created successfully",
		Code:    201,
	}
}

// GetUser retrieves a user by ID - WRONG IMPLEMENTATION
func (us *UserService) GetUser(ctx context.Context, userID string) Response {
	// WRONG: No ID validation
	// Should check if userID is empty, valid ObjectID format

	// WRONG: Not converting string to ObjectID
	filter := map[string]interface{}{"_id": userID} // Should be primitive.ObjectIDFromHex

	var user User
	err := us.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error", // WRONG: Should distinguish between "not found" and other errors
			Code:    500,
		}
	}

	return Response{
		Success: true,
		Data:    user,
		Message: "User retrieved successfully",
		Code:    200,
	}
}

// UpdateUser updates an existing user - WRONG IMPLEMENTATION
func (us *UserService) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) Response {
	// WRONG: No validation of userID or update fields
	// Should validate email format if provided, age range, etc.

	// WRONG: Not using $set operator
	update := req // Should be bson.M{"$set": req}

	// WRONG: Using wrong filter format
	filter := userID // Should be bson.M{"_id": objectID}

	result, err := us.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	// WRONG: Not checking if user was actually found and updated
	if result.ModifiedCount == 0 {
		// Should return "User not found" error
	}

	return Response{
		Success: true,
		Message: "User updated successfully",
		Code:    200,
	}
}

// DeleteUser deletes a user by ID - WRONG IMPLEMENTATION
func (us *UserService) DeleteUser(ctx context.Context, userID string) Response {
	// WRONG: No ID validation

	// WRONG: Using string instead of ObjectID
	filter := map[string]interface{}{"name": userID} // WRONG FIELD!

	result, err := us.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}

	// WRONG: Not checking if user was actually found and deleted
	if result.DeletedCount == 0 {
		// Should return "User not found" error
	}

	return Response{
		Success: true,
		Message: "User deleted successfully",
		Code:    200,
	}
}

// ListUsers retrieves all users - WRONG IMPLEMENTATION
func (us *UserService) ListUsers(ctx context.Context) Response {
	// WRONG: Using wrong filter format
	cursor, err := us.Collection.Find(ctx, "all") // Should be bson.M{} or bson.D{}
	if err != nil {
		return Response{
			Success: false,
			Error:   "Database error",
			Code:    500,
		}
	}
	defer cursor.Close(ctx)

	var users []User
	// WRONG: Not handling cursor.All error
	cursor.All(ctx, &users)

	return Response{
		Success: true,
		Data:    users,
		Message: "Users retrieved successfully",
		Code:    200,
	}
}

// ConnectMongoDB establishes connection to MongoDB - WRONG IMPLEMENTATION
func ConnectMongoDB(uri string) (*mongo.Client, error) {
	// WRONG: No URI validation
	// WRONG: Not actually implementing connection
	return nil, fmt.Errorf("ConnectMongoDB not implemented")
}
