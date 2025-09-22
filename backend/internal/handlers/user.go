package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"github.com/macielcr7/map-area-factions/backend/internal/repository"
	"github.com/macielcr7/map-area-factions/backend/internal/utils"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin moderator collaborator citizen"`
}

type UpdateUserRequest struct {
	Name   string `json:"name" validate:"omitempty,min=2,max=100"`
	Email  string `json:"email" validate:"omitempty,email"`
	Role   string `json:"role" validate:"omitempty,oneof=admin moderator collaborator citizen"`
	Active *bool  `json:"active" validate:"omitempty"`
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetUsers returns paginated list of users
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")
	role := c.Query("role", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	users, total, err := h.userRepo.GetUsers(page, limit, search, role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "USER_001",
			"message": "Failed to fetch users",
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
		"pagination": fiber.Map{
			"page":     page,
			"limit":    limit,
			"total":    total,
			"pages":    (total + limit - 1) / limit,
			"has_next": page*limit < total,
			"has_prev": page > 1,
		},
	})
}

// GetUser returns a specific user by ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "USER_002",
			"message": "Invalid user ID format",
		})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "USER_003",
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "USER_004",
			"message": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "USER_005",
			"message": "Invalid input data",
			"details": err,
		})
	}

	// Check if user already exists
	existingUser, _ := h.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return c.Status(409).JSON(fiber.Map{
			"error":   "Conflict",
			"code":    "USER_006",
			"message": "User with this email already exists",
		})
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Will be hashed in repository
		Role:     req.Role,
		Active:   true,
	}

	if err := h.userRepo.Create(user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "USER_007",
			"message": "Failed to create user",
		})
	}

	// Don't return password
	user.Password = ""

	return c.Status(201).JSON(fiber.Map{
		"message": "User created successfully",
		"user":    user,
	})
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "USER_008",
			"message": "Invalid user ID format",
		})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "USER_009",
			"message": "Invalid request body",
		})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Validation Error",
			"code":    "USER_010",
			"message": "Invalid input data",
			"details": err,
		})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "USER_011",
			"message": "User not found",
		})
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email already exists
		existingUser, _ := h.userRepo.GetByEmail(req.Email)
		if existingUser != nil && existingUser.ID != user.ID {
			return c.Status(409).JSON(fiber.Map{
				"error":   "Conflict",
				"code":    "USER_012",
				"message": "User with this email already exists",
			})
		}
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Active != nil {
		user.Active = *req.Active
	}

	if err := h.userRepo.Update(user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "USER_013",
			"message": "Failed to update user",
		})
	}

	// Don't return password
	user.Password = ""

	return c.JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    user,
	})
}

// DeleteUser soft deletes a user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "USER_014",
			"message": "Invalid user ID format",
		})
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   "Not Found",
			"code":    "USER_015",
			"message": "User not found",
		})
	}

	// Prevent deletion of current user
	currentUserID := c.Locals("user_id").(uuid.UUID)
	if user.ID == currentUserID {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"code":    "USER_016",
			"message": "Cannot delete your own user account",
		})
	}

	if err := h.userRepo.DeleteUser(user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"code":    "USER_017",
			"message": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
