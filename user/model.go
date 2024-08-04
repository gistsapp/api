package user

import (
	"database/sql"
	"errors"

	"github.com/gistapp/api/storage"
	"github.com/gofiber/fiber/v2/log"
)

type UserSQL struct {
	ID      sql.NullString
	Email   sql.NullString
	Name    sql.NullString
	Picture sql.NullString
}

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type UserModel interface {
	Save() (*User, error)
}

func (u *UserSQL) Save() (*User, error) {
	if u.Picture.Valid {
		u.Picture.String = "https://vercel.com/api/www/avatar/?u=" + u.Email.String + "&s=80"
		u.Picture.Valid = true
	}
	row, err := storage.Database.Query("INSERT INTO users(email, name, picture) VALUES ($1, $2, $3) RETURNING user_id, email, name, picture", u.Email.String, u.Name.String, u.Picture)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't create user")
	}
	var user User
	row.Next()
	err = row.Scan(&user.ID, &user.Email, &user.Name, &user.Picture)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find user")
	}

	return &user, nil
}
