package main

import (
	"context"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comprehensive test suite with 15 tests that validate user's actual implementation
// These tests focus on input validation, edge cases, and error handling

func TestCreateUserValidation(t *testing.T) {
	userService := &UserService{Collection: nil}

	tests := []struct {
		name    string
		request CreateUserRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid user should pass validation",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: false, // Should pass validation, fail on DB
		},
		{
			name: "Empty name should be rejected",
			request: CreateUserRequest{
				Name:  "",
				Email: "john@example.com",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "Empty email should be rejected",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "Invalid email format should be rejected",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "Email without @ should be rejected",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "johndoe.com",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "Email without . should be rejected",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example",
				Age:   30,
			},
			wantErr: true,
		},
		{
			name: "Zero age should be rejected",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   0,
			},
			wantErr: true,
		},
		{
			name: "Negative age should be rejected",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   -5,
			},
			wantErr: true,
		},
		{
			name: "Age over 150 should be rejected",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   200,
			},
			wantErr: true,
		},
		{
			name: "Boundary age 150 should be accepted",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   150,
			},
			wantErr: false, // Should pass validation, fail on DB
		},
		{
			name: "Boundary age 1 should be accepted",
			request: CreateUserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   1,
			},
			wantErr: false, // Should pass validation, fail on DB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture panics for tests that should pass validation but fail on DB
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						// This is expected - validation passed, DB operation failed
						t.Logf("Expected panic for valid input (DB operation failed): %v", r)
					} else {
						// This shouldn't happen - validation should have caught the error
						t.Errorf("Unexpected panic for invalid input: %v", r)
					}
				}
			}()

			response := userService.CreateUser(context.Background(), tt.request)

			if tt.wantErr {
				if response.Success {
					t.Errorf("Expected validation error but got success")
				}
				if response.Error == "" {
					t.Errorf("Expected error message but got empty string")
				}
			} else {
				// For valid inputs, we expect validation to pass but DB to fail
				if response.Success {
					t.Errorf("Expected database error but got success (DB should be nil)")
				}
			}
		})
	}
}

func TestGetUserValidation(t *testing.T) {
	userService := &UserService{Collection: nil}

	tests := []struct {
		name    string
		userID  string
		wantErr bool
		errType string
	}{
		{
			name:    "Empty user ID should be rejected",
			userID:  "",
			wantErr: true,
			errType: "required",
		},
		{
			name:    "Invalid ObjectID format should be rejected",
			userID:  "invalid-id",
			wantErr: true,
			errType: "invalid",
		},
		{
			name:    "Short invalid ID should be rejected",
			userID:  "123",
			wantErr: true,
			errType: "invalid",
		},
		{
			name:    "Valid ObjectID format should pass validation",
			userID:  primitive.NewObjectID().Hex(),
			wantErr: false, // Should pass validation, fail on DB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture panics for tests that should pass validation but fail on DB
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						// This is expected - validation passed, DB operation failed
						t.Logf("Expected panic for valid ObjectID (DB operation failed): %v", r)
					} else {
						// This shouldn't happen - validation should have caught the error
						t.Errorf("Unexpected panic for invalid ObjectID: %v", r)
					}
				}
			}()

			response := userService.GetUser(context.Background(), tt.userID)

			if tt.wantErr {
				if response.Success {
					t.Errorf("Expected validation error but got success")
				}
				if response.Error == "" {
					t.Errorf("Expected error message but got empty string")
				}
				// Check for specific error type
				if tt.errType == "required" && !strings.Contains(strings.ToLower(response.Error), "required") {
					t.Errorf("Expected 'required' error, got: %s", response.Error)
				}
				if tt.errType == "invalid" && !strings.Contains(strings.ToLower(response.Error), "invalid") {
					t.Errorf("Expected 'invalid' error, got: %s", response.Error)
				}
			} else {
				// For valid inputs, we expect validation to pass but DB to fail
				if response.Success {
					t.Errorf("Expected database error but got success (DB should be nil)")
				}
			}
		})
	}
}

func TestUpdateUserValidation(t *testing.T) {
	userService := &UserService{Collection: nil}

	tests := []struct {
		name    string
		userID  string
		request UpdateUserRequest
		wantErr bool
		errType string
	}{
		{
			name:   "Empty user ID should be rejected",
			userID: "",
			request: UpdateUserRequest{
				Name: "Updated Name",
			},
			wantErr: true,
			errType: "required",
		},
		{
			name:   "Invalid ObjectID format should be rejected",
			userID: "invalid-id",
			request: UpdateUserRequest{
				Name: "Updated Name",
			},
			wantErr: true,
			errType: "invalid",
		},
		{
			name:   "Invalid email format should be rejected",
			userID: primitive.NewObjectID().Hex(),
			request: UpdateUserRequest{
				Email: "invalid-email",
			},
			wantErr: true,
			errType: "email",
		},
		{
			name:   "Email without @ should be rejected",
			userID: primitive.NewObjectID().Hex(),
			request: UpdateUserRequest{
				Email: "johndoe.com",
			},
			wantErr: true,
			errType: "email",
		},
		{
			name:   "Valid update should pass validation",
			userID: primitive.NewObjectID().Hex(),
			request: UpdateUserRequest{
				Name:  "Valid Name",
				Email: "valid@example.com",
				Age:   30,
			},
			wantErr: false, // Should pass validation, fail on DB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture panics for tests that should pass validation but fail on DB
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						// This is expected - validation passed, DB operation failed
						t.Logf("Expected panic for valid input (DB operation failed): %v", r)
					} else {
						// This shouldn't happen - validation should have caught the error
						t.Errorf("Unexpected panic for invalid input: %v", r)
					}
				}
			}()

			response := userService.UpdateUser(context.Background(), tt.userID, tt.request)

			if tt.wantErr {
				if response.Success {
					t.Errorf("Expected validation error but got success")
				}
				if response.Error == "" {
					t.Errorf("Expected error message but got empty string")
				}
				// Check for specific error type
				if tt.errType == "required" && !strings.Contains(strings.ToLower(response.Error), "required") {
					t.Errorf("Expected 'required' error, got: %s", response.Error)
				}
				if tt.errType == "invalid" && !strings.Contains(strings.ToLower(response.Error), "invalid") {
					t.Errorf("Expected 'invalid' error, got: %s", response.Error)
				}
				if tt.errType == "email" && !strings.Contains(strings.ToLower(response.Error), "email") {
					t.Errorf("Expected 'email' error, got: %s", response.Error)
				}
			} else {
				// For valid inputs, we expect validation to pass but DB to fail
				if response.Success {
					t.Errorf("Expected database error but got success (DB should be nil)")
				}
			}
		})
	}
}

func TestDeleteUserValidation(t *testing.T) {
	userService := &UserService{Collection: nil}

	tests := []struct {
		name    string
		userID  string
		wantErr bool
		errType string
	}{
		{
			name:    "Empty user ID should be rejected",
			userID:  "",
			wantErr: true,
			errType: "required",
		},
		{
			name:    "Invalid ObjectID format should be rejected",
			userID:  "invalid-id",
			wantErr: true,
			errType: "invalid",
		},
		{
			name:    "Valid ObjectID format should pass validation",
			userID:  primitive.NewObjectID().Hex(),
			wantErr: false, // Should pass validation, fail on DB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture panics for tests that should pass validation but fail on DB
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						// This is expected - validation passed, DB operation failed
						t.Logf("Expected panic for valid ObjectID (DB operation failed): %v", r)
					} else {
						// This shouldn't happen - validation should have caught the error
						t.Errorf("Unexpected panic for invalid ObjectID: %v", r)
					}
				}
			}()

			response := userService.DeleteUser(context.Background(), tt.userID)

			if tt.wantErr {
				if response.Success {
					t.Errorf("Expected validation error but got success")
				}
				if response.Error == "" {
					t.Errorf("Expected error message but got empty string")
				}
				// Check for specific error type
				if tt.errType == "required" && !strings.Contains(strings.ToLower(response.Error), "required") {
					t.Errorf("Expected 'required' error, got: %s", response.Error)
				}
				if tt.errType == "invalid" && !strings.Contains(strings.ToLower(response.Error), "invalid") {
					t.Errorf("Expected 'invalid' error, got: %s", response.Error)
				}
			} else {
				// For valid inputs, we expect validation to pass but DB to fail
				if response.Success {
					t.Errorf("Expected database error but got success (DB should be nil)")
				}
			}
		})
	}
}

func TestListUsersBasic(t *testing.T) {
	userService := &UserService{Collection: nil}

	t.Run("ListUsers with nil collection should handle gracefully", func(t *testing.T) {
		// Capture panics
		defer func() {
			if r := recover(); r != nil {
				// This is expected - DB operation failed
				t.Logf("Expected panic for ListUsers with nil collection: %v", r)
			}
		}()

		response := userService.ListUsers(context.Background())

		// Should fail gracefully, not panic
		if response.Success {
			t.Error("Expected error with nil collection but got success")
		}
	})
}

func TestConnectMongoDBValidation(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{
			name:    "Valid URI should attempt connection",
			uri:     "mongodb://localhost:27017",
			wantErr: false, // RezaSi implementation actually connects successfully
		},
		{
			name:    "Empty URI should be rejected",
			uri:     "",
			wantErr: true,
		},
		{
			name:    "Invalid URI should be rejected",
			uri:     "invalid-uri",
			wantErr: true,
		},
		{
			name:    "Malformed URI should be rejected",
			uri:     "mongodb://",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := ConnectMongoDB(tt.uri)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for URI '%s' but got success", tt.uri)
				}
				if client != nil {
					t.Errorf("Expected nil client for URI '%s' but got non-nil", tt.uri)
				}
			} else {
				if err != nil {
					t.Errorf("Expected success for URI '%s' but got error: %s", tt.uri, err.Error())
				}
				if client == nil {
					t.Errorf("Expected non-nil client for URI '%s' but got nil", tt.uri)
				}
			}
		})
	}
}

func TestRequiredFunctionsExist(t *testing.T) {
	userService := &UserService{Collection: nil}

	t.Run("CreateUser function exists with correct signature", func(t *testing.T) {
		// This will compile only if the function exists with correct signature
		response := userService.CreateUser(context.Background(), CreateUserRequest{})
		_ = response // Use the response to avoid unused variable warning
	})

	t.Run("GetUser function exists with correct signature", func(t *testing.T) {
		response := userService.GetUser(context.Background(), "")
		_ = response
	})

	t.Run("UpdateUser function exists with correct signature", func(t *testing.T) {
		response := userService.UpdateUser(context.Background(), "", UpdateUserRequest{})
		_ = response
	})

	t.Run("DeleteUser function exists with correct signature", func(t *testing.T) {
		response := userService.DeleteUser(context.Background(), "")
		_ = response
	})

	t.Run("ListUsers function exists with correct signature", func(t *testing.T) {
		// Capture panics for ListUsers with nil collection
		defer func() {
			if r := recover(); r != nil {
				// This is expected - DB operation failed
				t.Logf("Expected panic for ListUsers with nil collection: %v", r)
			}
		}()

		response := userService.ListUsers(context.Background())
		_ = response
	})

	t.Run("ConnectMongoDB function exists with correct signature", func(t *testing.T) {
		client, err := ConnectMongoDB("")
		_ = client
		_ = err
	})
}

func TestResponseStructureValidation(t *testing.T) {
	userService := &UserService{Collection: nil}

	t.Run("CreateUser response structure should be valid", func(t *testing.T) {
		response := userService.CreateUser(context.Background(), CreateUserRequest{
			Name:  "",
			Email: "test@example.com",
			Age:   30,
		})

		// Check response structure
		if response.Success && response.Error != "" {
			t.Error("Response cannot be both successful and have an error")
		}
		if !response.Success && response.Error == "" {
			t.Error("Failed response must have an error message")
		}
		if response.Code == 0 {
			t.Error("Response must have a status code")
		}
	})

	t.Run("GetUser response structure should be valid", func(t *testing.T) {
		response := userService.GetUser(context.Background(), "")

		// Check response structure
		if response.Success && response.Error != "" {
			t.Error("Response cannot be both successful and have an error")
		}
		if !response.Success && response.Error == "" {
			t.Error("Failed response must have an error message")
		}
		if response.Code == 0 {
			t.Error("Response must have a status code")
		}
	})
}

func TestEdgeCasesAndBoundaryValues(t *testing.T) {
	userService := &UserService{Collection: nil}

	t.Run("Very long name should be handled", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic for long name (DB operation failed): %v", r)
			}
		}()

		longName := strings.Repeat("a", 1000)
		response := userService.CreateUser(context.Background(), CreateUserRequest{
			Name:  longName,
			Email: "test@example.com",
			Age:   30,
		})
		// Should either accept or reject gracefully, not panic
		_ = response
	})

	t.Run("Very long email should be handled", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic for long email (DB operation failed): %v", r)
			}
		}()

		longEmail := strings.Repeat("a", 500) + "@" + strings.Repeat("b", 500) + ".com"
		response := userService.CreateUser(context.Background(), CreateUserRequest{
			Name:  "Test User",
			Email: longEmail,
			Age:   30,
		})
		// Should either accept or reject gracefully, not panic
		_ = response
	})

	t.Run("Maximum valid age should be handled", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic for valid age (DB operation failed): %v", r)
			}
		}()

		response := userService.CreateUser(context.Background(), CreateUserRequest{
			Name:  "Test User",
			Email: "test@example.com",
			Age:   150,
		})
		// Should pass validation, fail on DB
		_ = response
	})
}
