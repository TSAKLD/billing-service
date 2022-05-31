package entity

import (
	"errors"
	"unicode/utf8"
)

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

func (u User) Validate() error {
	l := utf8.RuneCountInString(u.Name)
	if l < 3 || l > 50 {
		return errors.New("name must be longer than 3 symbols and longer than 50 symbols")
	}

	return nil
}
