package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"web_api/auth"
	"web_api/models"
	"web_api/responses"
)

// Rename these structs later
type User struct {
	us *models.UserService
}

type UserNew struct {
	Name		string	`json:"name"`
	Email 		string	`json:"email"`
	Password	string	`json:"password`
}

func NewUser(us *models.UserService) *User{
	return &User{
		us: us,
	}
}

// Create User
//
// POST /login
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var up UserNew
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&up)
	if err != nil {
		responses.JsonErrorResponse(w, err)
		return
	}

	// todo: add bcrypt + pepper
	nu := models.User{
		Name: up.Name,
		Email: up.Email,
		Password: up.Password,
	}

	lastID, err := u.us.Create(&nu)
	if err != nil {
		msg := "unable to complete create user request: " + err.Error()
		responses.SendErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}

	responses.SendDataResponse(w, http.StatusCreated, strconv.FormatInt(lastID, 10))
	return
}

type UserLogin struct {
	Email		string	`json:"email"`
	Password	string	`json:"password"`
}

// Login User
// Creates a jwt token
// POST /login
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	var ul UserLogin

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&ul)
	if err != nil {
		responses.JsonErrorResponse(w, err)
		return
	}

	// Check the user's credentials
	result, err := u.us.ByEmail(ul.Email)
	if err != nil {
		fmt.Println(err)
		msg := "login information does not match our records"
		responses.SendErrorResponse(w, http.StatusUnauthorized, msg)
		return
	}
	if result.Password != ul.Password {
		msg := "login information does not match our records"
		responses.SendErrorResponse(w, http.StatusUnauthorized, msg)
		return
	}

	userAuth, err := u.us.CreateAuth(result.ID)
	if err != nil {
		log.Println(err.Error())
		msg := http.StatusText(http.StatusInternalServerError)
		responses.SendErrorResponse(w, http.StatusInternalServerError, msg)
		return
	}

	userAuthorization := auth.UserAuth{}
	userAuthorization.UserUUID = userAuth.UserUUID
	userAuthorization.UserID = userAuth.UserID
	token, err := auth.CreateToken(userAuthorization)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	responses.SendDataResponse(w, http.StatusOK, token)
	return
}

// Logout User
// Logs the user out by deleting claims from the authorization table
// POST /logout
func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	token, err  := auth.ValidateJwtToken(bearerToken)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	ua, err := auth.RetrieveUserAuthorization(token)
	if err != nil {
		responses.SendErrorResponse(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err = u.us.DeleteAuth(ua)
	if err != nil {
		log.Println(err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	responses.SendDataResponse(w, http.StatusOK, nil)
	return
}

type UserUpdate struct {
	UserID 	    int		`json:"id"`
	Name		string	`json:"name"`
	Email 		string	`json:"email"`
}

// UPDATE User
// Updates the user's name, email address
// PUT /users/{id}
func (u *User) Update(w http.ResponseWriter, r *http.Request) {
	var upd models.User
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err = dec.Decode(&upd)
	if err != nil {
		responses.JsonErrorResponse(w, err)
		return
	}

	u.us.Update(userID, upd)
	responses.SendDataResponse(w, http.StatusOK, nil)
	return
}

// Delete User
// Deletes a User from the database
// POST /users/{id}
func (u *User) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		responses.SendErrorResponse(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	err = u.us.Delete(userID)
	if err != nil {
		msg := "unable to complete delete user request: " + err.Error()
		responses.SendErrorResponse(w, http.StatusBadRequest, msg)
		return
	}

	responses.SendDataResponse(w, http.StatusOK, nil)
	return
}