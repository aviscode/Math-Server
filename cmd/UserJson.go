package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	usersDb  *usersJson
	FileName = "users.json"
	once     = &sync.Once{}
)

type usersJson struct {
	Users map[string]*user `json:"usersJson"`
}

func NewUsersJson() *usersJson {
	once.Do(func() {
		usersDb = &usersJson{
			Users: make(map[string]*user, 32),
		}
		if err := usersDb.LoadUsers(); err != nil {
			panic(fmt.Sprintf("failed to load users.json file %v", err))
		}
	})
	return usersDb
}

func (U *usersJson) AddUser(user *user) error {
	if _, err := U.GetUser(user.UserName); err != ErrNoSuchUser {
		return ErrUserAlreadyExist
	}
	U.Users[user.UserName] = user
	return nil
}

func (U *usersJson) GetUser(userName string) (*user, error) {
	if v, ok := U.Users[userName]; !ok {
		return &user{}, ErrNoSuchUser
	} else {
		return v, nil
	}
}

func (U *usersJson) LoadUsers() error {
	dataFile, err := os.ReadFile(FileName)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(dataFile, &U.Users); err != nil {
		return err
	}
	return nil
}

func (U *usersJson) SaveUsers() error {
	data, err := json.Marshal(U.Users)
	if err != nil {
		return err
	}
	if err = os.WriteFile(FileName, data, 0666); err != nil {
		return err
	}
	return nil
}
