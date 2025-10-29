package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rd2w/rd2w-log/internal/shared/auth"
	"github.com/rd2w/rd2w-log/internal/shared/logger"
)

type jwtService struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
	logger          logger.Logger
}

// NewTokenService creates a new JWT token service
func NewTokenService(secret string, accessTTL, refreshTTL time.Duration, issuer string, logger logger.Logger) auth.TokenService {
	return &jwtService{
		secret:          secret,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
		issuer:          issuer,
		logger:          logger,
	}
}

type claims struct {
	UserID    string `json:"user_id"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func (s *jwtService) GenerateAccessToken(userID string) (string, error) {
	return s.generateToken(userID, s.accessTokenTTL, "access")
}

func (s *jwtService) GenerateRefreshToken(userID string) (string, error) {
	return s.generateToken(userID, s.refreshTokenTTL, "refresh")
}

func (s *jwtService) generateToken(userID string, ttl time.Duration, tokenType string) (string, error) {
	claims := &claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   userID,
			ID:        fmt.Sprintf("%s-%d", tokenType, time.Now().UnixNano()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *jwtService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		s.logger.Debug("Token validation failed", "error", err)
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		// Проверяем тип токена
		if claims.TokenType != "access" && claims.TokenType != "refresh" {
			return "", fmt.Errorf("invalid token type: %s", claims.TokenType)
		}
		return claims.UserID, nil
	}

	return "", fmt.Errorf("invalid token claims")
}

func (s *jwtService) ParseToken(tokenString string) (map[string]interface{}, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &claims{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*claims); ok {
		return map[string]interface{}{
			"user_id":    claims.UserID,
			"token_type": claims.TokenType,
			"exp":        claims.ExpiresAt,
			"iat":        claims.IssuedAt,
			"iss":        claims.Issuer,
		}, nil
	}

	return nil, fmt.Errorf("failed to parse token claims")
}

func (s *jwtService) RevokeToken(ctx context.Context, token string) error {
	// В будущем можно добавить blacklist в Redis
	s.logger.Debug("Token revocation requested", "token", token[:10]+"...")
	return nil
}

func (s *jwtService) IsTokenRevoked(ctx context.Context, token string) (bool, error) {
	// В будущем можно проверять blacklist в Redis
	return false, nil
}
