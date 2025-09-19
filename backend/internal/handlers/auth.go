package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/auth"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
	"github.com/macielcr7/map-area-factions/backend/internal/utils"
)

type AuthHandler struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         *models.User `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func NewAuthHandler(userRepo repository.UserRepository, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Invalid request body",
			Code:    "VALIDATION_001",
			Message: "Request body must be valid JSON",
		})
	}

	// Validate input
	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Validation failed",
			Code:    "VALIDATION_002",
			Message: err.Error(),
		})
	}

	// Find user
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse{
			Error:   "Invalid credentials",
			Code:    "AUTH_001",
			Message: "Email or password is incorrect",
		})
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(utils.ErrorResponse{
			Error:   "Invalid credentials",
			Code:    "AUTH_001",
			Message: "Email or password is incorrect",
		})
	}

	// Check if user is active
	if user.Status != "active" {
		return c.Status(401).JSON(utils.ErrorResponse{
			Error:   "Account inactive",
			Code:    "AUTH_007",
			Message: "Your account has been deactivated",
		})
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse{
			Error:   "Failed to generate token",
			Code:    "AUTH_008",
			Message: "Internal server error",
		})
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse{
			Error:   "Failed to generate refresh token",
			Code:    "AUTH_009",
			Message: "Internal server error",
		})
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24 hours
		User:         user,
	}

	return c.JSON(response)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(utils.ErrorResponse{
			Error:   "Invalid request body",
			Code:    "VALIDATION_001",
			Message: "Request body must be valid JSON",
		})
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse{
			Error:   "Invalid refresh token",
			Code:    "AUTH_010",
			Message: "Refresh token is invalid or expired",
		})
	}

	// Get user
	user, err := h.userRepo.GetByID(claims.UserID)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse{
			Error:   "User not found",
			Code:    "AUTH_011",
			Message: "User associated with token not found",
		})
	}

	// Generate new access token
	accessToken, err := h.jwtService.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return c.Status(500).JSON(utils.ErrorResponse{
			Error:   "Failed to generate token",
			Code:    "AUTH_008",
			Message: "Internal server error",
		})
	}

	response := LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   86400, // 24 hours
		User:        user,
	}

	return c.JSON(response)
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(401).JSON(utils.ErrorResponse{
			Error:   "User not authenticated",
			Code:    "AUTH_004",
			Message: "You must be logged in to access this resource",
		})
	}

	user, err := h.userRepo.GetByID(userID.(uuid.UUID))
	if err != nil {
		return c.Status(404).JSON(utils.ErrorResponse{
			Error:   "User not found",
			Code:    "USER_001",
			Message: "User profile not found",
		})
	}

	return c.JSON(user)
}