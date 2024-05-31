package models

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided was invalid")
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
	Name  string
	Email string `gorm:"not null;uniqueIndex"`
}

func (us *UserService) Create(user *User) error {
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
		{Name: "gg", Email: "gg@gg.com"},
		{Name: "mary", Email: "mm@mm.com"},
		{Name: "jane", Email: "jj@jj.com"},
		{Name: "nick", Email: "n@n.com"},
		{Name: "henry", Email: "hh@hh.com"},
		{Name: "bush", Email: "bb@bb.com"},
	}
	for _, user := range users {
		if err := us.db.Create(user).Error; err != nil {
			panic(err)
		}
	}
}
