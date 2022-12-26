package server

import (
	"MathServer/cmd"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name              string
		userName          string
		password          string
		x                 float64
		y                 float64
		LoginResponseCode int
		addResponseCode   int
		expectedResult    string
	}{
		{"login existing user but not in group A or B", "admin", "admin", 3.0, 3.0, 200, 200, "3.00 + 3.00 = 6.00"},
		{"login existing user but in group A", "adminA", "adminA", 3.0, 3.0, 200, 401, ""},
		{"login non existing user", "adminC", "adminC", 3.0, 3.0, 401, 401, ""},
		{"login with no body", "", "", 3.0, 3.0, 400, 401, ""},
	}
	cmd.FileName = "../users.json"
	InitServer("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Credentials{
				Password: tt.userName,
				Username: tt.password,
			}
			body, _ := json.Marshal(c)
			if tt.userName == "" {
				body = nil
			}
			r, err := http.NewRequest("GET", "#", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("faled to create request  error = %v,", err)
				return
			}
			w := httptest.NewRecorder()
			login(w, r)
			if w.Code != tt.LoginResponseCode {
				t.Errorf("login() = %d, want %d", w.Code, tt.LoginResponseCode)
				return
			}

			if tt.LoginResponseCode != 200 {
				return
			}
			cookie := w.Header().Get("Set-Cookie")
			FullSessionToken := strings.Split(cookie, ";")[0]
			sessionToken := strings.Split(FullSessionToken, "=")[1]

			p := params{
				X: tt.x,
				Y: tt.y,
			}
			body, _ = json.Marshal(p)
			if tt.userName == "" {
				body = nil
			}
			r, err = http.NewRequest("GET", "#", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("faled to create request  error = %v,", err)
				return
			}
			w = httptest.NewRecorder()
			r.AddCookie(&http.Cookie{
				Name:  SessionTokenName,
				Value: sessionToken,
			})
			add(w, r)
			if w.Code != tt.addResponseCode {
				t.Errorf("multiply() = %d, want %d", w.Code, tt.addResponseCode)
				return
			}
			if w.Body != nil && w.Body.String() != tt.expectedResult {
				t.Errorf("multiply() = %s, want %s", w.Body.String(), tt.expectedResult)
			}
		})
	}
}

func TestExponent(t *testing.T) {
	tests := []struct {
		name                 string
		userName             string
		password             string
		x                    float64
		y                    float64
		LoginResponseCode    int
		exponentResponseCode int
		expectedResult       string
	}{
		{"login existing user but not in group A or B", "admin", "admin", 3.0, 3.0, 200, 401, ""},
		{"login existing user but in group A", "adminA", "adminA", 3.0, 3.0, 200, 401, ""},
		{"login existing user but in group B", "adminB", "adminB", 3.0, 3.0, 200, 200, "3.00 ^ 3.00 = 27.00"},
		{"login non existing user", "adminC", "adminC", 3.0, 3.0, 401, 401, ""},
		{"login with no body", "", "", 3.0, 3.0, 400, 401, ""},
	}
	cmd.FileName = "../users.json"
	InitServer("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Credentials{
				Password: tt.userName,
				Username: tt.password,
			}
			body, _ := json.Marshal(c)
			if tt.userName == "" {
				body = nil
			}
			r, err := http.NewRequest("GET", "#", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("faled to create request  error = %v,", err)
				return
			}
			w := httptest.NewRecorder()
			login(w, r)
			if w.Code != tt.LoginResponseCode {
				t.Errorf("login() = %d, want %d", w.Code, tt.LoginResponseCode)
				return
			}

			if tt.LoginResponseCode != 200 {
				return
			}
			cookie := w.Header().Get("Set-Cookie")
			FullSessionToken := strings.Split(cookie, ";")[0]
			sessionToken := strings.Split(FullSessionToken, "=")[1]

			p := params{
				X: tt.x,
				Y: tt.y,
			}
			body, _ = json.Marshal(p)
			if tt.userName == "" {
				body = nil
			}
			r, err = http.NewRequest("GET", "#", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("faled to create request  error = %v,", err)
				return
			}
			w = httptest.NewRecorder()
			r.AddCookie(&http.Cookie{
				Name:  SessionTokenName,
				Value: sessionToken,
			})
			exponent(w, r)
			if w.Code != tt.exponentResponseCode {
				t.Errorf("multiply() = %d, want %d", w.Code, tt.exponentResponseCode)
				return
			}
			if w.Body != nil && w.Body.String() != tt.expectedResult {
				t.Errorf("multiply() = %s, want %s", w.Body.String(), tt.expectedResult)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name         string
		userName     string
		password     string
		responseCode int
	}{
		{"login existing user", "admin", "admin", 200},
		{"login non existing user", "adminC", "adminC", 401},
		{"login with no body", "", "", 400},
	}
	cmd.FileName = "../users.json"
	InitServer("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Credentials{
				Password: tt.userName,
				Username: tt.password,
			}
			body, _ := json.Marshal(c)
			if tt.userName == "" {
				body = nil
			}
			r, err := http.NewRequest("GET", "#", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("faled to create request  error = %v,", err)
				return
			}
			w := httptest.ResponseRecorder{}
			login(&w, r)
			if w.Code != tt.responseCode {
				t.Errorf("login() = %d, want %d", w.Code, tt.responseCode)
				return
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name              string
		userName          string
		password          string
		x                 float64
		y                 float64
		LoginResponseCode int
		multiResponseCode int
		expectedResulte   string
	}{
		{"login existing user but not in group A", "admin", "admin", 3.0, 3.0, 200, 401, ""},
		{"login existing user but in group A", "adminA", "adminA", 3.0, 3.0, 200, 200, "3.00 * 3.00 = 9.00"},
		{"login non existing user", "adminC", "adminC", 3.0, 3.0, 401, 401, ""},
		{"login with no body", "", "", 3.0, 3.0, 400, 401, ""},
	}
	InitServer("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Credentials{
				Password: tt.userName,
				Username: tt.password,
			}
			body, _ := json.Marshal(c)
			if tt.userName == "" {
				body = nil
			}
			r, err := http.NewRequest("GET", "#", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("faled to create request  error = %v,", err)
				return
			}
			w := httptest.NewRecorder()
			login(w, r)
			if w.Code != tt.LoginResponseCode {
				t.Errorf("login() = %d, want %d", w.Code, tt.LoginResponseCode)
				return
			}

			if tt.LoginResponseCode != 200 {
				return
			}
			cookie := w.Header().Get("Set-Cookie")
			FullSessionToken := strings.Split(cookie, ";")[0]
			sessionToken := strings.Split(FullSessionToken, "=")[1]

			p := params{
				X: tt.x,
				Y: tt.y,
			}
			body, _ = json.Marshal(p)
			if tt.userName == "" {
				body = nil
			}
			r, err = http.NewRequest("GET", "#", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("faled to create request  error = %v,", err)
				return
			}
			w = httptest.NewRecorder()
			r.AddCookie(&http.Cookie{
				Name:  SessionTokenName,
				Value: sessionToken,
			})
			multiply(w, r)
			if w.Code != tt.multiResponseCode {
				t.Errorf("multiply() = %d, want %d", w.Code, tt.multiResponseCode)
				return
			}
			if w.Body != nil && w.Body.String() != tt.expectedResulte {
				t.Errorf("multiply() = %s, want %s", w.Body.String(), tt.expectedResulte)
			}
		})
	}
}
