package router

import (
	"crudapi/controllers"
	"crudapi/repository"
	"crudapi/service"
	"database/sql"

	"github.com/gorilla/mux"
)

func InitRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	r.HandleFunc("/users", userController.GetUsers).Methods("GET")

	return r
}
