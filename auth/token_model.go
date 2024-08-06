package auth

import (
	"database/sql"
	"errors"

	"github.com/gistapp/api/storage"
	"github.com/gofiber/fiber/v2/log"
)

type TokenType string

const (
	LocalAuth TokenType = "local_auth"
)

type TokenSQL struct {
	ID      sql.NullString
	Type    sql.NullString
	Value   sql.NullString
	Keyword sql.NullString
}

type Token struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	Keyword string `json:"keyword"`
}

func (t *TokenSQL) Save() (*Token, error) {
	row, err := storage.Database.Query("INSERT INTO token(type, value, keyword) VALUES ($1, $2, $3) RETURNING token_id, type, value, keyword", t.Type.String, t.Value.String, t.Keyword.String)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't create token")
	}
	var token Token
	row.Next()
	err = row.Scan(&token.ID, &token.Type, &token.Value, &token.Keyword)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find token")
	}
	return &token, nil
}

func (t *TokenSQL) Get() (*Token, error) {
	row, err := storage.Database.Query("SELECT token_id, type, value, keyword FROM token WHERE type = $1 AND keyword = $2", t.Type.String, t.Keyword.String)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find token")
	}
	var token Token
	row.Next()
	err = row.Scan(&token.ID, &token.Type, &token.Value, &token.Keyword)
	if err != nil {
		log.Error(err)
		return nil, errors.New("couldn't find token")
	}
	return &token, nil
}

func (t *Token) Delete() error {
	_, err := storage.Database.Query("DELETE FROM token WHERE value = $1 and keyword = $2 and type = $3", t.Value, t.Keyword, t.Type)
	if err != nil {
		log.Error(err)
		return errors.New("couldn't delete token")
	}
	return nil
}
