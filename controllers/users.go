package controllers

import (
	"fmt"
	"net/http"

	"github.com/Georgeygigz/go-commuters-hub/models"
	"github.com/Georgeygigz/go-commuters-hub/views"
)

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

// POST login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm

	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(w, "Invalid email address")
	case models.ErrInvalidPassword:
		fmt.Fprintln(w, "Invalid password provided.")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

// New is used to render the form where a user can
// create a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm

	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userData := fmt.Sprintf("%+v", user)
	fmt.Fprintln(w, userData)
}
