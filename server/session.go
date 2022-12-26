package server

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"sync"
	"time"
)

var (
	SessionTokenName  = "session_token"
	ErrNoSuchSession  = errors.New("no such session")
	ErrSessionExpired = errors.New("session expired")
	once              sync.Once
	s                 *sessions
)

type session struct {
	username string
	expiry   time.Time
}

// we'll use this method later to determine if the session has expired
func (s *session) isExpired() bool {
	if s == nil {
		return true
	}
	return s.expiry.Before(time.Now())
}

type sessions struct {
	sessionsMap map[string]*session
	mu          sync.Mutex
}

func Sessions() *sessions {
	once.Do(func() {
		s = &sessions{
			mu:          sync.Mutex{},
			sessionsMap: make(map[string]*session, 32),
		}
	})
	return s
}

func (s *sessions) CreateNewSession(userName string) (string, time.Time) {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessionsMap[sessionToken] = &session{
		username: userName,
		expiry:   expiresAt,
	}
	return sessionToken, expiresAt
}

func (s *sessions) verifySessionTokenAndGetUserName(r *http.Request) (string, error) {
	c, err := r.Cookie(SessionTokenName)
	if err != nil {
		return "", err
	}
	userSession, exists := s.sessionsMap[c.Value]
	if !exists {
		return "", ErrNoSuchSession
	}
	if userSession.isExpired() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.sessionsMap, c.Value)
		return "", ErrSessionExpired
	}
	return userSession.username, nil
}
