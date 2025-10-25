package domain_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rd2w/rd2w-log/internal/auth/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewSession(t *testing.T) {
	userID := uuid.New()
	accessToken := "access_token"
	refreshToken := "refresh_token"
	expiresAt := time.Now().Add(24 * time.Hour)

	session := domain.NewSession(userID, accessToken, refreshToken, expiresAt)

	assert.NotNil(t, session)
	assert.NotEmpty(t, session.ID)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, accessToken, session.AccessToken)
	assert.Equal(t, refreshToken, session.RefreshToken)
	assert.Equal(t, expiresAt, session.ExpiresAt)
	assert.WithinDuration(t, time.Now(), session.CreatedAt, time.Second)
}

func TestSession_IsExpired(t *testing.T) {
	userID := uuid.New()

	t.Run("not expired", func(t *testing.T) {
		session := domain.NewSession(userID, "token", "refresh", time.Now().Add(time.Hour))
		assert.False(t, session.IsExpired())
	})

	t.Run("expired", func(t *testing.T) {
		session := domain.NewSession(userID, "token", "refresh", time.Now().Add(-time.Hour))
		assert.True(t, session.IsExpired())
	})

	t.Run("just expired", func(t *testing.T) {
		session := domain.NewSession(userID, "token", "refresh", time.Now().Add(-time.Second))
		assert.True(t, session.IsExpired())
	})
}

func TestSession_Refresh(t *testing.T) {
	userID := uuid.New()
	session := domain.NewSession(userID, "old_token", "refresh", time.Now().Add(time.Hour))

	oldExpiresAt := session.ExpiresAt
	newToken := "new_token"
	newExpiresAt := time.Now().Add(2 * time.Hour)

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	session.Refresh(newToken, newExpiresAt)

	assert.Equal(t, newToken, session.AccessToken)
	assert.Equal(t, newExpiresAt, session.ExpiresAt)
	assert.True(t, session.ExpiresAt.After(oldExpiresAt))
}
