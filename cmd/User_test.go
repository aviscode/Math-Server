package cmd

import "testing"

func TestGetUserGroup(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		password string
		group    Group
		wantErr  bool
	}{
		{"groupA", "avi", "avi", NoGroup, false},
		{"groupB", "avir", "avir", GroupB, false},
		{"groupC", "avirs", "avirs", GroupB, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser(tt.userName, tt.password, tt.group)
			if err != nil {
				t.Errorf("falied to create user %v: %v", tt.name, err)
				return
			}
			if got := u.GetUserGroup(); got != tt.group {
				t.Errorf("GetUserGroup() = %v, want %v", got, tt.group)
			}
		})
	}
}

func TestVerifyPass(t *testing.T) {
	tests := []struct {
		name          string
		userName      string
		password      string
		inputPassword string
		group         Group
		wantErr       bool
	}{
		{"good pass", "avi", "avi", "avi", NoGroup, false},
		{"good pass", "avir", "avir", "avir", NoGroup, false},
		{" bad pass", "avirs", "avirs", "avi", NoGroup, true},
		{" bad pass", "avirs", "avirs", "", NoGroup, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser(tt.userName, tt.password, tt.group)
			if err != nil {
				t.Errorf("falied to create user %v: %v", tt.name, err)
				return
			}
			if err := u.VerifyPassword(tt.inputPassword); (err != nil) != tt.wantErr {
				t.Errorf("VerifyPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
