package users

import (
	"errors"
	"fmt"

	"github.com/Rowkash/go-gin-auth/internal/common"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

type UserFilter struct {
	ID    *int    `form:"id"`
	Email *string `form:"email"`
	Name  *string `form:"name"`
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Create(data UserInsertInput) (*User, error) {
	user := User{
		Email:    data.Email,
		Name:     data.Name,
		Password: data.Password,
	}

	result := s.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *Service) FindOne(opts UserFilter) (*User, error) {
	var user User

	err := s.db.Scopes(s.filterUser(opts)).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.NewNotFoundError("User not found")
		}

		return nil, err
	}
	return &user, nil
}

func (s *Service) getPage(opts AdminUsersPageRequest) ([]User, error) {
	var users []User

	offset := (opts.Page - 1) * opts.Limit

	err := s.db.
		Scopes(s.filterUser(opts.UserFilter), s.sortUser(opts.PaginationQuery)).
		Limit(opts.Limit).
		Offset(offset).
		Find(&users).
		Error

	return users, err
}

func (s *Service) filterUser(opts UserFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if opts.ID != nil {
			db = db.Where("id = ?", *opts.ID)
		}
		if opts.Email != nil {
			db = db.Where("email = ?", *opts.Email)
		}
		if opts.Name != nil {
			db = db.Where("name LIKE ?", "%"+*opts.Name+"%")
		}
		return db
	}
}

func (s *Service) sortUser(opts PaginationQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		dbColumn := `"createdAt"`

		switch opts.SortBy {
		case SortByName:

			dbColumn = `"name"`
		case SortByEmail:
			dbColumn = `"email"`
		}

		return db.Order(fmt.Sprintf("%s %s", dbColumn, opts.OrderSort))
	}
}
