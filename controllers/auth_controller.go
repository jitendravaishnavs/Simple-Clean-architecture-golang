package controllers

import (
	"crudapi/model"
	"crudapi/utils"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (c *AuthController) Signup(w http.ResponseWriter, r *http.Request) {
	var u model.User
	json.NewDecoder(r.Body).Decode(&u)

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(hashedPass)

	err := c.AuthRepo.CreateUser(u)
	if err != nil {
		http.Error(w, "Error creating user", 500)
		return
	}

	token, _ := utils.GenerateToken(u.Email)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var creds model.User
	json.NewDecoder(r.Body).Decode(&creds)

	user, err := c.AuthRepo.FindByEmail(creds.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)) != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	token, _ := utils.GenerateToken(user.Email)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
