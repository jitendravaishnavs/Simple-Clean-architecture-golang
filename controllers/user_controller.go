package controllers

import (
	"crudapi/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(s service.UserService) *UserController {
	return &UserController{UserService: s}
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.UserService.GetUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func GetUsersWithSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := query.Get("search")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	// Default pagination
	page := 1
	limit := 10

	if pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	offset := (page - 1) * limit

	// Base query
	baseQuery := "SELECT id, name, email, image FROM users WHERE 1"
	var args []interface{}

	// If search present
	if search != "" {
		baseQuery += " AND (name LIKE ? OR email LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	baseQuery += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		http.Error(w, "Query error", 500)
		return
	}
	defer rows.Close()

	var users []map[string]string
	for rows.Next() {
		var id int
		var name, email, image string
		rows.Scan(&id, &name, &email, &image)

		users = append(users, map[string]string{
			"id":    fmt.Sprint(id),
			"name":  name,
			"email": email,
			"image": "/uploads/" + image,
		})
	}

	json.NewEncoder(w).Encode(users)
}
