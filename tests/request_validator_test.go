package e_tests

import (
	"errors"
	"testing"

	"elemental/plugins/request-validator"
)

// Mocking data for the test
var testData = []request_validator.User{
	{ID: 1, Name: "John Doe", Age: 25, IsActive: true},
	{ID: 2, Name: "Jane Smith", Age: 30, IsActive: true},
	{ID: 3, Name: "Bob Johnson", Age: 40, IsActive: false},
	{ID: 4, Name: "Alice Williams", Age: 22, IsActive: true},
	{ID: 5, Name: "Mike Davis", Age: 35, IsActive: true},
}

func TestValidateStructWithDB(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		input         request_validator.User // User struct to validate
		expectedError error                   // Expected validation error
	}{
		{
			name: "ValidUser",
			input: request_validator.User{
				ID:       6,
				Name:     "Emily Watson",
				Age:      20,
				IsActive: true,
			},
			expectedError: nil, // No validation error expected
		},
		{
			name: "DuplicateID",
			input: request_validator.User{
				ID:       1, // ID already exists in the test data
				Name:     "Duplicate User",
				Age:      25,
				IsActive: true,
			},
			expectedError: errors.New("ID already exists"),
		},
		{
			name: "InvalidAge",
			input: request_validator.User{
				ID:       7,
				Name:     "Invalid Age User",
				Age:      10, // Age is less than 18
				IsActive: true,
			},
			expectedError: errors.New("Age is not greater than 18"),
		},
		{
			name: "InvalidIsActive",
			input: request_validator.User{
				ID:       8,
				Name:     "Invalid IsActive User",
				Age:      20,
				IsActive: false, // IsActive should be true
			},
			expectedError: errors.New("IsActive is not true"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := request_validator.ValidateStructWithDB( tc.input)
			if (err == nil && tc.expectedError != nil) || (err != nil && tc.expectedError == nil) || (err != nil && tc.expectedError != nil && err.Error() != tc.expectedError.Error()) {
				t.Errorf("Unexpected error: %v, expected: %v", err, tc.expectedError)
			}
		})
	}
}
