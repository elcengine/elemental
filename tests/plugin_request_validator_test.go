package e_tests

import (
	// "errors"
	"testing"

	"github.com/elcengine/elemental/plugins/request-validator"
)

type UserValidation struct {
	ID       int    `validate:"exists=user_table,UserId"`
	Name     string `validate:"exists=user_table,Name"`
	Age      int    `validate:"IsGreater=20"`
	IsActive bool   `validate:"isTrue"`
}


func TestValidateStructWithDB(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		input         UserValidation // User struct to validate
		expectedError error                   // Expected validation error
	}{
		{
			name: "ValidUser",
			input: UserValidation{
				ID:       6,
				Name:     "Emily Watson",
				Age:      20,
				IsActive: true,
			},
			expectedError: nil, // No validation error expected
		},
		{
			name: "DuplicateID",
			input: UserValidation{
				ID:       1, 
				Age:      25,
				IsActive: true,
			},
			expectedError: nil,
		},
		{
			name: "InvalidAge",
			input: UserValidation{
				ID:       7,
				Name:     "Invalid Age User",
				Age:      10, 
				IsActive: true,
			},
			expectedError: nil,
		},
		{
			name: "InvalidIsActive",
			input: UserValidation{
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
