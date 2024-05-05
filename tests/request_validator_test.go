package e_tests

import (
	// "errors"
	"testing"

	"elemental/plugins/request-validator"
)


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
				ID:       1, 
				Age:      25,
				IsActive: true,
			},
			expectedError: nil,
		},
		{
			name: "InvalidAge",
			input: request_validator.User{
				ID:       7,
				Name:     "Invalid Age User",
				Age:      10, 
				IsActive: true,
			},
			expectedError: nil,
		},
		{
			name: "InvalidIsActive",
			input: request_validator.User{
				ID:       8,
				Name:     "Invalid IsActive User",
				Age:      20,
				IsActive: false, 
			},
			expectedError: nil,
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
