package model

// TokenPair содержит пару access и refresh токенов
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LoginCredentials содержит данные для входа
type LoginCredentials struct {
	Callsign string `json:"callsign"`
	Password string `json:"password"`
}

// RegisterRequest содержит данные для регистрации
type RegisterRequest struct {
	Callsign string `json:"callsign"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Timezone string `json:"timezone"`
}

// UpdateProfileRequest содержит данные для обновления профиля
type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Timezone string `json:"timezone"`
}
