package server

import (
	"net/http"
	"testing"
	"time"
)

func TestIsExpired(t *testing.T) {
	tests := []struct {
		name   string
		s      session
		expire bool
	}{
		{"session no expire", session{"aa", time.Now().Add(120 * time.Second)}, false},
		{"session  expired", session{"aa", time.Now()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.isExpired(); got != tt.expire {
				t.Errorf("isExpired() = %v, want %v", got, tt.expire)
			}
		})
	}
}

func TestVerifySessionTokenAndGetUserName(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		expire   bool
	}{
		{"session no expire", "aa", false},
		{"session expired", "bb", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessionToken, expiresAt := Sessions().CreateNewSession(tt.userName)
			r, err := http.NewRequest("GET", "#", nil)
			if err != nil {
				t.Errorf("faled to create request  error = %v,", err)
				return
			}
			r.AddCookie(&http.Cookie{
				Name:    SessionTokenName,
				Value:   sessionToken,
				Expires: expiresAt,
			})
			if tt.expire {
				Sessions().sessionsMap[sessionToken].expiry = time.Now()
			}
			got, err := s.verifySessionTokenAndGetUserName(r)
			if err != nil && !tt.expire {
				t.Errorf("verifySessionTokenAndGetUserName() error = %v,", err)
				return
			}
			if !tt.expire && got != tt.userName {
				t.Errorf("verifySessionTokenAndGetUserName() got = %v, want %v", got, tt.userName)
			}
		})
	}
}
