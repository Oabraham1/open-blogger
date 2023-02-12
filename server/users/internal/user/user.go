package user

import (
	"errors"
	"time"
)

type User struct {
	ID        string
	Username  string
	FirstName string
	LastName  string
	Email     string
	ImageURL  string
	Created   time.Time
	Updated   time.Time
}

func NewUser(id, username, firstName, lastName, email, imageURL string) *User {
	return &User{
		ID:        id,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		ImageURL:  imageURL,
		Created:   time.Now(),
		Updated:   time.Now(),
	}
}

func (user *User) UpdateUser(update *User) error {
	if user == nil || update == nil {
		return errors.New("user update failed")
	}
	user.Username = update.Username
	user.FirstName = update.FirstName
	user.LastName = update.LastName
	user.Email = update.Email
	user.ImageURL = update.ImageURL
	user.Updated = time.Now()
	return nil
}

func (user *User) Validate() error {
	if user == nil {
		return errors.New("user validation failed")
	}
	if user.Username == "" {
		return errors.New("user validation failed")
	}
	return nil
}
