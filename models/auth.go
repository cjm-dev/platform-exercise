package models

import (
	"fmt"
	"github.com/google/uuid"
	"web_api/auth"
)

type Auth struct {
	ID 			int		`db:"id"`
	UserID 		int 	`db:"user_id"`
	UserUUID	string  `db:"auth_uuid"`
}

func (us *UserService) CreateAuth(userID int)  (*Auth, error){
	var userAuthorization Auth
	userAuthorization.UserID = userID
	userAuthorization.UserUUID = uuid.New().String()

	fmt.Println("ID =", userAuthorization.UserID)
	// insert
	query := `INSERT INTO authentication (user_id, auth_uuid) VALUES (:user_id, :auth_uuid)`
	_, err := us.db.NamedExec(query, userAuthorization)
	if err != nil {
		return nil, err
	}

	return &userAuthorization, nil
}

//Once a user row in the auth table
func (us *UserService) DeleteAuth(ua *auth.UserAuth) error {
	userID := ua.UserID
	userUUID := ua.UserUUID

	fmt.Println("userid", userID)
	fmt.Println("uuid", userUUID)
	query := `DELETE FROM authentication WHERE user_id = ? and auth_uuid = ? LIMIT 1`
	_, err := us.db.Exec(query, userID, userUUID)
	if err != nil {
		fmt.Println("delete auth error")
		return err
	}

	return nil
}

