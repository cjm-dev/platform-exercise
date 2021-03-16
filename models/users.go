package models

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"
)

type UserService struct {
	db *sqlx.DB
}

type User struct {
	ID				int			`db:"id"`
	Name			string		`db:"name"`
	Email 			string		`db:"email"`
	Password		string		`db:"password"`
	CreatedAt		time.Time 	`db:"created_at"`
	UpdatedAt		time.Time	`db:"updated_at"`
}

func NewUserService() (*UserService, error) {
	userDBCredential := os.Getenv("DB_USER") + ":" + os.Getenv("PASSWORD")
	dbHost := "@tcp(" + os.Getenv("LOCALHOST") + ":" + os.Getenv("PORT") + ")/" + os.Getenv("DB_NAME")
	dataSource := userDBCredential + dbHost + "?parseTime=true"
	db, err := sqlx.Open(os.Getenv("DB_DRIVER"), dataSource)
	if err != nil {
		return nil, err
	}
	return &UserService{
		db: db,
	}, nil
}

func (us *UserService) Create(user *User) (int64, error) {
	query := `INSERT INTO users (name, email, password) VALUES (:name, :email, :password)`
	result , err := us.db.NamedExec(query, user)
	if err != nil {
		return 0, err
	}

	li, err := result.LastInsertId()
	fmt.Println("li", li)
	if err != nil {
		return 0, err
	}

	return li, nil
}

// todo errors
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	query := `SELECT * from users where email = ?`
	err := us.db.Get(&user, query, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

//todo errors that happen here?
func (us *UserService) Update (userID int, userUpdate User) error {
	// todo resource not found error
	fmt.Println("In update model")
	fmt.Printf("%#v\n", userUpdate)
	tx := us.db.MustBegin()
	if userUpdate.Name != "" {
		fmt.Println("in Name", userUpdate.Name, userID)
		tx.MustExec(tx.Rebind("UPDATE users SET name = ? WHERE id = ?"), userUpdate.Name, userID)
	}
	if userUpdate.Email != "" {
		tx.MustExec(tx.Rebind("UPDATE users SET email = ? WHERE id = ?"), userUpdate.Email, userID)
	}
	tx.Commit()

	return nil
}

func (us *UserService) Delete(userID int) error {
	query := `DELETE FROM users where id = ?`
	result, err := us.db.Exec(query, userID)
	if err == nil {
		count, err := result.RowsAffected()
		if count == 0 {
			err := errors.New("user id not found")
			return err
		}
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) Close() error {
	return us.db.Close()
}


