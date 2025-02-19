package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Register(u *User) (*User, error)
	Login(e string) (*User, error)
}

type UserMysqlRepository struct {
	db gorm.DB
}

func NewUserMysqlRepository(db gorm.DB) UserRepository {
	return &UserMysqlRepository{db: db}
}

func (r *UserMysqlRepository) Register(u *User) (*User, error) {
	var user_exists User
	if err := r.db.Where("email = ?", u.Email).First(&user_exists).Error; err == nil {
		return nil, gorm.ErrDuplicatedKey
	}

	if err := r.db.Save(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserMysqlRepository) Login(email string) (*User, error) {
	var user_exists User
	if err := r.db.Where("email = ?", email).First(&user_exists).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return &user_exists, nil
}
