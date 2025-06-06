package repository

import (
	"crudapi/model"
	"database/sql"
)

type UserRepo interface {
	GetAll() ([]model.User, error)
}

type userRepoImpl struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepoImpl{DB: db}
}

func (r *userRepoImpl) GetAll() ([]model.User, error) {
	rows, err := r.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}
	return users, nil
}
