package model_test

import (
	"testing"
	"time"

	"github.com/rd2w/rd2w-log/internal/auth/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		callsign    string
		email       string
		password    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid user",
			callsign: "RD2W",
			email:    "test@example.com",
			password: "securepassword123",
			wantErr:  false,
		},
		{
			name:        "empty callsign",
			callsign:    "",
			email:       "test@example.com",
			password:    "securepassword123",
			wantErr:     true,
			errContains: "invalid callsign",
		},
		{
			name:        "empty email",
			callsign:    "RD2W",
			email:       "",
			password:    "securepassword123",
			wantErr:     true,
			errContains: "invalid email",
		},
		{
			name:        "short password",
			callsign:    "RD2W",
			email:       "test@example.com",
			password:    "short",
			wantErr:     true,
			errContains: "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := model.NewUser(tt.callsign, tt.email, tt.password, "Test User", "Location", "UTC")

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, tt.callsign, user.Callsign)
			assert.Equal(t, tt.email, user.Email)
			assert.NotEmpty(t, user.Password)
			assert.NotEqual(t, tt.password, user.Password) // Password should be hashed
			assert.NotEmpty(t, user.ID)
			assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
			assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
		})
	}
}

func TestUser_CheckPassword(t *testing.T) {
	user, err := model.NewUser("RD2W", "test@example.com", "securepassword123", "Test User", "Location", "UTC")
	require.NoError(t, err)

	// Test correct password
	assert.True(t, user.CheckPassword("securepassword123"))

	// Test incorrect password
	assert.False(t, user.CheckPassword("wrongpassword"))

	// Test empty password
	assert.False(t, user.CheckPassword(""))
}

func TestUser_UpdateProfile(t *testing.T) {
	user, err := model.NewUser("RD2W", "test@example.com", "securepassword123", "Old Name", "Old Location", "UTC")
	require.NoError(t, err)

	oldUpdatedAt := user.UpdatedAt

	// Wait a moment to ensure time difference
	time.Sleep(time.Millisecond * 10)

	user.UpdateProfile("New Name", "New Location", "Europe/Moscow")

	assert.Equal(t, "New Name", user.Name)
	assert.Equal(t, "New Location", user.Location)
	assert.Equal(t, "Europe/Moscow", user.Timezone)
	assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
}

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
	}{
		{
			name: "valid user",
			user: &model.User{
				Callsign: "RD2W",
				Email:    "test@example.com",
				Password: "hashedpassword",
			},
			wantErr: false,
		},
		{
			name: "empty callsign",
			user: &model.User{
				Callsign: "",
				Email:    "test@example.com",
				Password: "hashedpassword",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			user: &model.User{
				Callsign: "RD2W",
				Email:    "",
				Password: "hashedpassword",
			},
			wantErr: true,
		},
		{
			name: "empty password",
			user: &model.User{
				Callsign: "RD2W",
				Email:    "test@example.com",
				Password: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
