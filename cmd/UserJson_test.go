package cmd

import (
	"os"
	"reflect"
	"testing"
)

func initTest() {
	FileName = "testUsers.json"
	users := usersJson{Users: map[string]*user{}}
	user, _ := NewUser("admin", "admin", NoGroup)
	users.AddUser(user)
	user, _ = NewUser("adminA", "adminA", GroupA)
	users.AddUser(user)
	user, _ = NewUser("adminB", "adminB", GroupB)
	users.AddUser(user)
	user, _ = NewUser("adminC", "adminC", NoGroup)
	users.AddUser(user)
	users.SaveUsers()
}

func TestNewUsersJson(t *testing.T) {
	initTest()
	defer os.RemoveAll(FileName)
	t.Run("test singleTone", func(t *testing.T) {
		users := NewUsersJson()
		if got := NewUsersJson(); !reflect.DeepEqual(got, users) {
			t.Errorf("NewUsersJson() = %v, want %v", got, users)
		}
	})
}

func TestAddUser(t *testing.T) {
	initTest()
	defer os.RemoveAll(FileName)
	users := NewUsersJson()
	type args struct {
		userName string
		password string
		group    Group
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{" new user no group", args{userName: "avi", password: "avi", group: NoGroup}, false},
		{" new user groupA", args{userName: "avir", password: "avisr", group: NoGroup}, false},
		{" new user groupb", args{userName: "avirs", password: "avirs", group: NoGroup}, false},
		{" add new user with existing user name ", args{userName: "avirs", password: "avirs", group: NoGroup}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.args.userName, tt.args.password, tt.args.group)
			if err != nil {
				t.Errorf("falied to create user %v: %v", tt.name, err)
			}
			if err := users.AddUser(user); (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	initTest()
	users := NewUsersJson()
	defer os.RemoveAll(FileName)
	tests := []struct {
		name     string
		userName string
		wantErr  bool
	}{
		{"get an existing user", "admin", false},
		{"get an existing user", "adminA", false},
		{"get an existing user", "adminB", false},
		{"get an existing user", "adminC", false},
		{"get an  non existing user", "admind", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := users.GetUser(tt.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.UserName != tt.userName {
				t.Errorf("GetUser() got user  = %s, want %s", got.UserName, tt.userName)
			}
		})
	}
}

func TestSaveAndLoadUsers(t *testing.T) {
	FileName = "testUsers.json"
	users := usersJson{Users: map[string]*user{}}
	user, _ := NewUser("admin", "admin", NoGroup)
	users.AddUser(user)
	user, _ = NewUser("adminA", "adminA", GroupA)
	users.AddUser(user)
	user, _ = NewUser("adminB", "adminB", GroupB)
	users.AddUser(user)
	user, _ = NewUser("adminC", "adminC", NoGroup)
	users.AddUser(user)
	users.SaveUsers()
	defer os.RemoveAll(FileName)
	users = *NewUsersJson()
	tests := []struct {
		name     string
		userName string
		wantErr  bool
	}{
		{"get an existing user", "admin", false},
		{"get an existing user", "adminA", false},
		{"get an existing user", "adminB", false},
		{"get an existing user", "adminC", false},
		{"get an non existing user", "admind", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := users.GetUser(tt.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.UserName != tt.userName {
				t.Errorf("GetUser() got user  = %s, want %s", got.UserName, tt.userName)
			}
		})
	}
}
