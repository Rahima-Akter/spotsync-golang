package service

import (
	"errors"
	"time"

	"github.com/Rahima-Akter/spotsync-golang/internal/config"
	"github.com/Rahima-Akter/spotsync-golang/internal/dto"
	"github.com/Rahima-Akter/spotsync-golang/internal/models"
	"github.com/Rahima-Akter/spotsync-golang/internal/repository"
	"github.com/Rahima-Akter/spotsync-golang/internal/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles all authentication business logic
// This is like combining passport.js + bcrypt + JWT in Express.js
type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Register user
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	// default role if not provided
	if req.Role == "" {
		req.Role = "driver"
	}

	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, utils.ErrDuplicateEmail
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create the user model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	// Save to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Return the response DTO (no password!)
	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Login user and returns a JWT token
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrInvalidCredentials
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Return the response with token and user data
	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	// Create the claims (payload) for the JWT
	// Like: jwt.sign({ userId: user.id, role: user.role }, secret, { expiresIn: '24h' })
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.config.JWTExpiryHours)).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	signedToken, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
