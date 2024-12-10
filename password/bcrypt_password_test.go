package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcryptPasswordAlgorithm_HidePassword(t *testing.T) {
	// Set up the test cases
	testCases := []struct {
		name     string
		password string
		salt     string
	}{
		{
			name:     "Test case 1",
			password: "password123",
			salt:     "salt123",
		},
		{
			name:     "Test case 2",
			password: "password456",
			salt:     "salt456",
		},
	}

	// Run the test cases
	for _, tc := range testCases {
		bpa := NewBcryptPasswordAlgorithm()

		t.Run(tc.name, func(t *testing.T) {
			// Create a new bcrypt password algorithm

			// Generate the hashed password
			hashedPassword, err := bpa.HidePassword(tc.password, tc.salt)
			assert.NoError(t, err)

			// Compare the hashed password with the expected result
			assert.NotEmpty(t, hashedPassword)

			valid := bpa.CheckPassword(tc.password, tc.salt, hashedPassword)
			assert.True(t, valid)
		})
	}
}
