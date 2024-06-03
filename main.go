package main

import (
	"fmt"
	"net/http"

	"github.com/Georgeygigz/go-commuters-hub/controllers"
	"github.com/Georgeygigz/go-commuters-hub/models"
	"github.com/gorilla/mux"
)

var (
	host     = "localhost"
	port     = 5432
	user     = "gigz"
	password = "2416gigz1994"
	name     = "commuters"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	// us.DestructiveReset()
	// us.AutoMigrate()
	// us.SeedUserData()
	user, err := us.ByID(1)
	// user, err := us.ByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	var h http.Handler = http.HandlerFunc(notFound)

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.NotFoundHandler = h

	http.ListenAndServe(":4000", r)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>PAGE NOT FOUND</h1>")
}

// func must(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
