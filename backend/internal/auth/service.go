package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Jerinji2016/halooid/backend/internal/models"
	"github.com/Jerinji2016/halooid/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Common errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrTokenBlacklisted   = errors.New("token has been blacklisted")
	ErrUserInactive       = errors.New("user is inactive")
)

// Config holds the configuration for the auth service
type Config struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

// Service provides authentication functionality
type Service interface {
	// Register registers a new user
	Register(ctx context.Context, reg models.UserRegistration) (*models.UserResponse, error)

	// Login authenticates a user and returns a token pair
	Login(ctx context.Context, login models.UserLogin) (*models.TokenPair, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, req models.RefreshTokenRequest) (*models.TokenPair, error)

	// ValidateToken validates a token and returns the claims
	ValidateToken(ctx context.Context, tokenString string, tokenType models.TokenType) (*models.TokenClaims, error)

	// Logout invalidates a token
	Logout(ctx context.Context, tokenString string) error

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error)
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
	config      Config
}

// NewService creates a new auth service
func NewService(userRepo repository.UserRepository, redisClient *redis.Client, config Config) Service {
	return &serviceImpl{
		userRepo:    userRepo,
		redisClient: redisClient,
		config:      config,
	}
}

// Register registers a new user
func (s *serviceImpl) Register(ctx context.Context, reg models.UserRegistration) (*models.UserResponse, error) {
	user, err := models.NewUser(reg)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// Login authenticates a user and returns a token pair
func (s *serviceImpl) Login(ctx context.Context, login models.UserLogin) (*models.TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if !user.CheckPassword(login.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate token pair
	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *serviceImpl) RefreshToken(ctx context.Context, req models.RefreshTokenRequest) (*models.TokenPair, error) {
	// Validate refresh token
	claims, err := s.ValidateToken(ctx, req.RefreshToken, models.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Blacklist the old refresh token
	err = s.blacklistToken(ctx, req.RefreshToken, claims.ExpiresAt.Sub(time.Now()))
	if err != nil {
		return nil, err
	}

	// Generate new token pair
	tokenPair, err := s.generateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return tokenPair, nil
}

// ValidateToken validates a token and returns the claims
func (s *serviceImpl) ValidateToken(ctx context.Context, tokenString string, tokenType models.TokenType) (*models.TokenClaims, error) {
	// Check if token is blacklisted
	blacklisted, err := s.isTokenBlacklisted(ctx, tokenString)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, ErrTokenBlacklisted
	}

	// Determine the secret based on token type
	var secret string
	if tokenType == models.AccessToken {
		secret = s.config.AccessTokenSecret
	} else {
		secret = s.config.RefreshTokenSecret
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Validate token type
	claimTokenType, ok := claims["token_type"].(string)
	if !ok || models.TokenType(claimTokenType) != tokenType {
		return nil, ErrInvalidToken
	}

	// Extract user ID
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Extract email
	email, ok := claims["email"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}

	// Extract issued at
	issuedAtFloat, ok := claims["iat"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	issuedAt := time.Unix(int64(issuedAtFloat), 0)

	// Extract expiration
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, ErrInvalidToken
	}
	expiresAt := time.Unix(int64(expFloat), 0)

	// Extract roles (if present)
	var roles []string
	if rolesInterface, ok := claims["roles"]; ok {
		if rolesArray, ok := rolesInterface.([]interface{}); ok {
			roles = make([]string, 0, len(rolesArray))
			for _, role := range rolesArray {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
	}

	// Extract permissions (if present)
	var permissions []string
	if permissionsInterface, ok := claims["permissions"]; ok {
		if permissionsArray, ok := permissionsInterface.([]interface{}); ok {
			permissions = make([]string, 0, len(permissionsArray))
			for _, permission := range permissionsArray {
				if permissionStr, ok := permission.(string); ok {
					permissions = append(permissions, permissionStr)
				}
			}
		}
	}

	// Extract organization ID (if present)
	var orgID uuid.UUID
	if orgIDStr, ok := claims["org_id"].(string); ok {
		if parsedOrgID, err := uuid.Parse(orgIDStr); err == nil {
			orgID = parsedOrgID
		}
	}

	return &models.TokenClaims{
		UserID:      userID,
		Email:       email,
		TokenType:   tokenType,
		Roles:       roles,
		Permissions: permissions,
		OrgID:       orgID,
		IssuedAt:    issuedAt,
		ExpiresAt:   expiresAt,
	}, nil
}

// Logout invalidates a token
func (s *serviceImpl) Logout(ctx context.Context, tokenString string) error {
	// Validate the token first
	claims, err := s.ValidateToken(ctx, tokenString, models.AccessToken)
	if err != nil {
		return err
	}

	// Blacklist the token
	return s.blacklistToken(ctx, tokenString, claims.ExpiresAt.Sub(time.Now()))
}

// GetUserByID retrieves a user by ID
func (s *serviceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}

// generateTokenPair generates an access token and refresh token pair
func (s *serviceImpl) generateTokenPair(user *models.User) (*models.TokenPair, error) {
	now := time.Now()
	accessTokenExpiry := now.Add(s.config.AccessTokenExpiry)
	refreshTokenExpiry := now.Add(s.config.RefreshTokenExpiry)

	// Get user roles for the default organization
	// In a real application, this would be determined based on the request
	defaultOrgID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Create a roleRepository to get user roles
	// This is not ideal as it creates a circular dependency, but it's a simple solution for now
	// In a real application, you would inject the roleRepository into the auth service
	db := s.userRepo.GetDB()
	roleRepo := repository.NewPostgresRoleRepository(db)

	// Get user roles
	roles, err := roleRepo.GetUserRoles(context.Background(), user.ID, defaultOrgID)
	if err != nil {
		// If there's an error getting roles, we'll just continue without them
		// This is not ideal, but it's better than failing the token generation
		roles = []models.Role{}
	}

	// Extract role names and permissions
	roleNames := make([]string, 0, len(roles))
	permissions := make([]string, 0)
	permissionMap := make(map[string]bool)

	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
		for _, permission := range role.Permissions {
			if !permissionMap[permission.Name] {
				permissions = append(permissions, permission.Name)
				permissionMap[permission.Name] = true
			}
		}
	}

	// Create access token with roles and permissions
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     user.ID.String(),
		"email":       user.Email,
		"token_type":  string(models.AccessToken),
		"roles":       roleNames,
		"permissions": permissions,
		"org_id":      defaultOrgID.String(),
		"iat":         now.Unix(),
		"exp":         accessTokenExpiry.Unix(),
		"iss":         s.config.Issuer,
	})

	// Sign access token
	accessTokenString, err := accessToken.SignedString([]byte(s.config.AccessTokenSecret))
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"token_type": string(models.RefreshToken),
		"iat":        now.Unix(),
		"exp":        refreshTokenExpiry.Unix(),
		"iss":        s.config.Issuer,
	})

	// Sign refresh token
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.RefreshTokenSecret))
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.config.AccessTokenExpiry.Seconds()),
	}, nil
}

// blacklistToken adds a token to the blacklist
func (s *serviceImpl) blacklistToken(ctx context.Context, tokenString string, expiration time.Duration) error {
	// Use Redis to store blacklisted tokens
	key := "blacklist:" + tokenString
	return s.redisClient.Set(ctx, key, 1, expiration).Err()
}

// isTokenBlacklisted checks if a token is blacklisted
func (s *serviceImpl) isTokenBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	key := "blacklist:" + tokenString
	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
