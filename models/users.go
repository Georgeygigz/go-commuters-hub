package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound        = errors.New("models: resource not found")
	ErrInvalidID       = errors.New("models: ID provided was invalid")
	ErrInvalidPassword = errors.New("models: incorrect password provided")
	userPwPepper       = "secret-random-string"
)

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open(postgres.Open(connectionInfo), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &UserService{
		db: db,
	}, nil

}

type UserService struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;uniqueIndex"`
	Password     string `gorm:"_"`
	PasswordHash string `gorm:"not null"`
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	pwBytes := []byte(password + userPwPepper)
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), pwBytes)
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}

}

func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pwBytes), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

func (us *UserService) DestructiveReset() error {
	migrator := us.db.Migrator()
	if migrator.HasTable(&User{}) {
		return migrator.DropTable(&User{})
	}
	us.AutoMigrate()
	// us.seedUserData()
	return nil
}

func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (us *UserService) SeedUserData() {
	users := []*User{
		{Name: "gg", Email: "gg@gg.com", PasswordHash: "hash-1"},
		{Name: "mary", Email: "mm@mm.com", PasswordHash: "hash-2"},
		{Name: "jane", Email: "jj@jj.com", PasswordHash: "hash-3"},
		{Name: "nick", Email: "n@n.com", PasswordHash: "hash-4"},
		{Name: "henry", Email: "hh@hh.com", PasswordHash: "hash-5"},
		{Name: "bush", Email: "bb@bb.com", PasswordHash: "hash-6"},
	}
	for _, user := range users {
		if err := us.db.Create(user).Error; err != nil {
			panic(err)
		}
	}
}
