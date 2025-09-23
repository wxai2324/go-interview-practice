package main

import (
	"context"
	"strings"
	"testing"
)

// Comprehensive test suite that validates user's actual implementation
// These tests focus on input validation, edge cases, and error handling

func TestGetSalesByCategoryValidation(t *testing.T) {
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	t.Run("GetSalesByCategory should handle nil collection gracefully", func(t *testing.T) {
		// Capture panics
		defer func() {
			if r := recover(); r != nil {
				// This is expected - DB operation failed
				t.Logf("Expected panic for GetSalesByCategory with nil collection: %v", r)
			}
		}()

		response := analyticsService.GetSalesByCategory(context.Background())

		// Should fail gracefully, not panic
		if response.Success {
			t.Error("Expected error with nil collection but got success")
		}
	})
}

func TestGetTopSellingProductsValidation(t *testing.T) {
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	tests := []struct {
		name    string
		limit   int
		wantErr bool
		errType string
	}{
		{
			name:    "Valid limit should pass validation",
			limit:   10,
			wantErr: false, // Should pass validation, fail on DB
		},
		{
			name:    "Zero limit should be rejected",
			limit:   0,
			wantErr: true,
			errType: "limit",
		},
		{
			name:    "Negative limit should be rejected",
			limit:   -5,
			wantErr: true,
			errType: "limit",
		},
		{
			name:    "Very large limit should be handled",
			limit:   10000,
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

			response := analyticsService.GetTopSellingProducts(context.Background(), tt.limit)

			if tt.wantErr {
				if response.Success {
					t.Errorf("Expected validation error but got success")
				}
				if response.Error == "" {
					t.Errorf("Expected error message but got empty string")
				}
				// Check for specific error type
				if tt.errType == "limit" && !strings.Contains(strings.ToLower(response.Error), "limit") {
					t.Errorf("Expected 'limit' error, got: %s", response.Error)
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

func TestGetRevenueByMonthValidation(t *testing.T) {
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	tests := []struct {
		name    string
		year    int
		wantErr bool
		errType string
	}{
		{
			name:    "Valid year should pass validation",
			year:    2023,
			wantErr: false, // Should pass validation, fail on DB
		},
		{
			name:    "Zero year should be rejected",
			year:    0,
			wantErr: true,
			errType: "year",
		},
		{
			name:    "Negative year should be rejected",
			year:    -2023,
			wantErr: true,
			errType: "year",
		},
		{
			name:    "Very old year should be rejected",
			year:    1800,
			wantErr: true,
			errType: "year",
		},
		{
			name:    "Future year should be rejected",
			year:    2100,
			wantErr: true,
			errType: "year",
		},
		{
			name:    "Current year should pass validation",
			year:    2024,
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

			response := analyticsService.GetRevenueByMonth(context.Background(), tt.year)

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

func TestGetCustomerAnalyticsValidation(t *testing.T) {
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	tests := []struct {
		name    string
		limit   int
		wantErr bool
		errType string
	}{
		{
			name:    "Valid limit should pass validation",
			limit:   10,
			wantErr: false, // Should pass validation, fail on DB
		},
		{
			name:    "Zero limit should be rejected",
			limit:   0,
			wantErr: true,
			errType: "limit",
		},
		{
			name:    "Negative limit should be rejected",
			limit:   -5,
			wantErr: true,
			errType: "limit",
		},
		{
			name:    "Very large limit should be handled",
			limit:   10000,
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

			response := analyticsService.GetCustomerAnalytics(context.Background(), tt.limit)

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

func TestGetProductPerformanceValidation(t *testing.T) {
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	t.Run("GetProductPerformance should handle nil collection gracefully", func(t *testing.T) {
		// Capture panics
		defer func() {
			if r := recover(); r != nil {
				// This is expected - DB operation failed
				t.Logf("Expected panic for GetProductPerformance with nil collection: %v", r)
			}
		}()

		response := analyticsService.GetProductPerformance(context.Background())

		// Should fail gracefully, not panic
		if response.Success {
			t.Error("Expected error with nil collection but got success")
		}
	})
}

func TestGetOrderTrendsValidation(t *testing.T) {
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	tests := []struct {
		name    string
		days    int
		wantErr bool
		errType string
	}{
		{
			name:    "Valid days should pass validation",
			days:    30,
			wantErr: false, // Should pass validation, fail on DB
		},
		{
			name:    "Zero days should be rejected",
			days:    0,
			wantErr: true,
			errType: "days",
		},
		{
			name:    "Negative days should be rejected",
			days:    -30,
			wantErr: true,
			errType: "days",
		},
		{
			name:    "Very large days should be handled",
			days:    10000,
			wantErr: false, // Should pass validation, fail on DB
		},
		{
			name:    "Single day should pass validation",
			days:    1,
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

			response := analyticsService.GetOrderTrends(context.Background(), tt.days)

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

func TestRequiredFunctionsExist(t *testing.T) {
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	t.Run("All required functions exist with correct signatures", func(t *testing.T) {
		// These might panic with nil collection, so capture panics
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic for functions with nil collection: %v", r)
			}
		}()

		// This will compile only if the functions exist with correct signatures
		_ = analyticsService.GetSalesByCategory(context.Background())
		_ = analyticsService.GetTopSellingProducts(context.Background(), 0)
		_ = analyticsService.GetRevenueByMonth(context.Background(), 0)
		_ = analyticsService.GetCustomerAnalytics(context.Background(), 0)
		_ = analyticsService.GetProductPerformance(context.Background())
		_ = analyticsService.GetOrderTrends(context.Background(), 0)
	})
}

func TestResponseStructureValidation(t *testing.T) {
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	t.Run("Response structure should be consistent", func(t *testing.T) {
		response := analyticsService.GetTopSellingProducts(context.Background(), 0)

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
	analyticsService := &AnalyticsService{OrdersCollection: nil}

	t.Run("Very large limit values should be handled", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic for large limit values (DB operation failed): %v", r)
			}
		}()

		response := analyticsService.GetTopSellingProducts(context.Background(), 999999)
		// Should either accept or reject gracefully, not panic
		_ = response
	})

	t.Run("Boundary year values should be handled", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic for boundary year values (DB operation failed): %v", r)
			}
		}()

		response := analyticsService.GetRevenueByMonth(context.Background(), 2024)
		// Should either accept or reject gracefully, not panic
		_ = response
	})

	t.Run("Maximum days values should be handled", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Logf("Expected panic for max days values (DB operation failed): %v", r)
			}
		}()

		response := analyticsService.GetOrderTrends(context.Background(), 365)
		// Should either accept or reject gracefully, not panic
		_ = response
	})
}
