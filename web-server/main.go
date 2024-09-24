package main

import (
	"errors"
	"fmt"
	"net/http"
	"web-server/handlers"
)

const PORT = ":3000"

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handlers.Root)
	mux.HandleFunc("POST /users", handlers.CreateUser)
	mux.HandleFunc("GET /users", handlers.GetUsers)
	mux.HandleFunc("GET /users/{id}", handlers.GetUser)
	mux.HandleFunc("DELETE /users/{id}", handlers.DeleteUser)

	fmt.Println("Server listening on", PORT)
	err := http.ListenAndServe(PORT, mux)
	if err != nil {
		panic(errors.New("Server panic"))
	}
}
