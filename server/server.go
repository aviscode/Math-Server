package server

import (
	"MathServer/cmd"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

var (
	users cmd.IUsers
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type params struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//verify the user exist
	UserInfo, err := users.GetUser(creds.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//verify the user's password
	if err := UserInfo.VerifyPassword(creds.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken, expiresAt := Sessions().CreateNewSession(creds.Username)
	http.SetCookie(w, &http.Cookie{
		Name:    SessionTokenName,
		Value:   sessionToken,
		Expires: expiresAt,
	})
	w.WriteHeader(http.StatusOK)
}

func add(w http.ResponseWriter, r *http.Request) {
	params, err := extractPramsFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userName, err := Sessions().verifySessionTokenAndGetUserName(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userInfo, err := users.GetUser(userName); err != nil || userInfo.GetUserGroup() != cmd.NoGroup {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("%.2f + %.2f = %.2f", params.X, params.Y, params.X+params.Y)))
}

func multiply(w http.ResponseWriter, r *http.Request) {
	params, err := extractPramsFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userName, err := Sessions().verifySessionTokenAndGetUserName(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userInfo, err := users.GetUser(userName); err != nil || userInfo.GetUserGroup() != cmd.GroupA {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("%.2f * %.2f = %.2f", params.X, params.Y, params.X*params.Y)))
}

func exponent(w http.ResponseWriter, r *http.Request) {
	params, err := extractPramsFromBody(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userName, err := Sessions().verifySessionTokenAndGetUserName(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userInfo, err := users.GetUser(userName); err != nil || userInfo.GetUserGroup() != cmd.GroupB {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("%.2f ^ %.2f = %.2f", params.X, params.Y, math.Pow(params.X, params.Y))))
}

func extractPramsFromBody(r *http.Request) (*params, error) {
	var params params
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return nil, err
	}
	return &params, nil
}

func InitServer(port string) *http.Server {
	users = cmd.NewUsersJson()
	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/add", add)
	mux.HandleFunc("/multiply", multiply)
	mux.HandleFunc("/exponent", exponent)
	return &http.Server{
		Handler: mux,
		Addr:    port,
	}
}
