package cmd

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Group int8

const (
	NoGroup Group = iota
	GroupA
	GroupB
	ErrGroup
)

var (
	ErrNoSuchUser         = errors.New("no such user")
	ErrPasswordDosNotMuch = errors.New("the password dos not much")
	ErrUserAlreadyExist   = errors.New("user already exist")
)

type user struct {
	UserName       string `json:"UserName"`
	HashedPassword []byte `json:"HashedPassword"`
	Group          Group  `json:"Group"`
}

func (u *user) GetUserGroup() Group {
	if u == nil {
		return ErrGroup
	}
	return u.Group
}

func NewUser(userName, password string, group Group) (*user, error) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return &user{UserName: userName, HashedPassword: hashPass, Group: group}, nil
}

// VerifyPassword Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return error on failure.
func (u *user) VerifyPassword(inputPassword string) error {
	if u == nil {
		return ErrNoSuchUser
	}
	if err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(inputPassword)); err != nil {
		return ErrPasswordDosNotMuch
	}
	return nil
}
