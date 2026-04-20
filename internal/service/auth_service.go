package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eulerbutcooler/virtus/internal/config"
	"github.com/eulerbutcooler/virtus/internal/domain"
	redisrepo "github.com/eulerbutcooler/virtus/internal/repository/redis"
	"github.com/eulerbutcooler/virtus/pkg/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Returned on every successful auth operation
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// JWT payload
type Claims struct {
	UserID uuid.UUID       `json:"uid"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type RegisterInput struct {
	Email    string
	Name     string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthService struct {
	users  domain.UserRepository
	cache  *redisrepo.Cache
	jwtCfg config.JWTConfig
}

func NewAuthService(users domain.UserRepository, cache *redisrepo.Cache, jwtCfg config.JWTConfig) *AuthService {
	return &AuthService{users: users, cache: cache, jwtCfg: jwtCfg}
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*domain.User, *TokenPair, error) {
	_, err := s.users.GetByEmail(ctx, in.Email)
	if err == nil {
		return nil, nil, domain.ErrConflict
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, nil, fmt.Errorf("auth.Register: email lookup: %w", err)
	}
	hash, err := crypto.HashPassword(in.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("auth.Register: hash password: %w", err)
	}
	user, err := s.users.Create(ctx, domain.CreateUserParams{
		Email:    in.Email,
		Name:     in.Name,
		Password: hash,
		Role:     domain.RoleMember,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("auth.Register: create user: %w", err)
	}
	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("auth.Register: issue tokens: %w", err)
	}
	return user, pair, nil
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*domain.User, *TokenPair, error) {
	user, err := s.users.GetByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, nil, domain.ErrUnauthorized
		}
		return nil, nil, fmt.Errorf("auth.Login: lookup: %w", err)
	}
	if err := crypto.CheckPassword(user.PasswordHash, in.Password); err != nil {
		return nil, nil, domain.ErrUnauthorized
	}
	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return nil, nil, fmt.Errorf("auth.Login: issue tokens: %w", err)
	}
	return user, pair, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	rtKey := rtCacheKey(claims.UserID, claims.ID)
	exists, err := s.cache.Exists(ctx, rtKey)
	if err != nil || !exists {
		return nil, domain.ErrUnauthorized
	}
	user, err := s.users.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: get user: %w", err)
	}
	_ = s.cache.Del(ctx, rtKey)
	pair, err := s.issueTokenPair(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("auth.Refresh: issue tokens: %w", err)
	}
	return pair, nil
}

// Invalidates the refresh token in Redis
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return nil
	}
	_ = s.cache.Del(ctx, rtCacheKey(claims.UserID, claims.ID))
	return nil
}

// Will be used by middleware to validate the access token
func (s *AuthService) ValidateAccessToken(tokenStr string) (*Claims, error) {
	claims, err := s.parseToken(tokenStr)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return claims, nil
}

func (s *AuthService) issueTokenPair(ctx context.Context, user *domain.User) (*TokenPair, error) {
	accessExpiry := time.Now().Add(s.jwtCfg.AccessTokenTTL)
	accessTok, err := s.mintToken(user, accessExpiry, "access")
	if err != nil {
		return nil, fmt.Errorf("mint access token: %w", err)
	}
	refreshExpiry := time.Now().Add(s.jwtCfg.RefreshTokenTTL)
	refreshTok, err := s.mintToken(user, refreshExpiry, "refresh")
	if err != nil {
		return nil, fmt.Errorf("mint refresh token: %w", err)
	}

	refreshClaims, err := s.parseToken(refreshTok)
	if err != nil {
		return nil, fmt.Errorf("parse refresh claims: %w", err)
	}
	rtKey := rtCacheKey(user.ID, refreshClaims.ID)
	if err := s.cache.Set(ctx, rtKey, true, s.jwtCfg.RefreshTokenTTL); err != nil {
		return nil, fmt.Errorf("store refresh token in cache: %w", err)
	}
	return &TokenPair{
		AccessToken:  accessTok,
		RefreshToken: refreshTok,
		ExpiresAt:    accessExpiry,
	}, nil
}

// Signs a new JWT for the user with a unique jti.
// Audience claim differentiates between refresh and access tokens so they can't be substituted for each other.
func (s *AuthService) mintToken(user *domain.User, expiry time.Time, tokenType string) (string, error) {
	jti, err := crypto.GenerateToken(16)
	if err != nil {
		return "", fmt.Errorf("generate jti: %w", err)
	}
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiry),
			Audience:  jwt.ClaimStrings{"virtus:" + tokenType},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtCfg.Secret))
}

// Returns parsed claims on success.
// Will add audience enforcing after http handlers.
func (s *AuthService) parseToken(tokenStr string) (*Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(s.jwtCfg.Secret), nil
		},
	)
	if err != nil || !token.Valid {
		return nil, domain.ErrUnauthorized
	}
	return &claims, nil
}

// Returns a redis namespaced key which allows per-session invalidation
func rtCacheKey(userID uuid.UUID, jti string) string {
	return fmt.Sprintf("rt:%s:%s", userID, jti)
}
