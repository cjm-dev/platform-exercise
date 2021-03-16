package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"web_api/controllers"
	"web_api/models"
	"github.com/joho/godotenv"
)

func  main() {

	// Load environment variables
	err := godotenv.Load(".env")
	// Setup db connection
	us, err := models.NewUserService()
	if err != nil {
		panic(err)
	}
	defer us.Close()

	userController := controllers.NewUser(us)
	r := mux.NewRouter()

	r.HandleFunc("/signup", userController.Create).Methods("POST")
	r.HandleFunc("/login", userController.Login).Methods("POST")
	r.HandleFunc("/logout", userController.Logout).Methods("POST")
	r.HandleFunc("/users/{id}", userController.Delete).Methods("DELETE")
	r.HandleFunc("/users/{id}", userController.Update).Methods("PUT")

	http.ListenAndServe(":3000", r)
}
