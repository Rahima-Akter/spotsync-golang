package repository

import (
	"errors"

	"github.com/Rahima-Akter/spotsync-golang/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles all database operations for users
// This is like creating a Prisma service:
//
//	const userService = {
//	  findByEmail: (email) => prisma.user.findUnique({ where: { email } }),
//	  create: (data) => prisma.user.create({ data }),
//	}

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database
// Like: prisma.user.create({ data: user })
func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		// Check if the error is due to duplicate email
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return errors.New("email already exists")
		}
		return result.Error
	}
	return nil
}

// FindByEmail looks up a user by their email address
// Like: prisma.user.findUnique({ where: { email } })
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		// If no user found, return nil without error
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// FindByID looks up a user by their ID
// Like: prisma.user.findUnique({ where: { id } })
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
