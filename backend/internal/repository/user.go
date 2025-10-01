package repository

import (
	"github.com/google/uuid"
	"github.com/macielcr7/map-area-factions/backend/internal/auth"
	"github.com/macielcr7/map-area-factions/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	DeleteUser(user *models.User) error
	List(limit, offset int, role string) ([]*models.User, int64, error)
	GetUsers(page, limit int, search, role string) ([]models.User, int, error)
	Count() (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	// Hash password before saving
	if user.Password != "" {
		hashedPassword, err := auth.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) DeleteUser(user *models.User) error {
	return r.db.Delete(user).Error
}

func (r *userRepository) List(limit, offset int, role string) ([]*models.User, int64, error) {
	var users []*models.User
	var count int64

	query := r.db.Model(&models.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *userRepository) GetUsers(page, limit int, search, role string) ([]models.User, int, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	// Apply search filter
	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Apply role filter
	if role != "" {
		query = query.Where("role = ?", role)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}