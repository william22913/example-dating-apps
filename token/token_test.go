package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	jwtTokenValidator := NewJWTTokenValidator("token_key", time.Duration(1*time.Hour))

	token, err := jwtTokenValidator.GenerateToken(1, "123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := jwtTokenValidator.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), payload.UserID)
	assert.Equal(t, "user_token", payload.StandardClaims.Subject)
	assert.Equal(t, "user", payload.StandardClaims.Audience)
	assert.Equal(t, "dating-apps", payload.StandardClaims.Issuer)
	assert.NotEmpty(t, payload.StandardClaims.ExpiresAt)
	assert.NotEmpty(t, payload.StandardClaims.IssuedAt)
	assert.Equal(t, "123", payload.StandardClaims.Id)

}
