package cmd

type IUsers interface {
	AddUser(user *user) error
	GetUser(userName string) (*user, error)
	LoadUsers() error
	SaveUsers() error
}
