package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error on reading body", http.StatusBadRequest)
		return
	}

	user := new(User)
	err = json.Unmarshal(body, user)
	// err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error on parsing user data", http.StatusBadRequest)
		return
	}

	if len(user.Name) < 3 {
		http.Error(w, "Please provide user.name at least in 3 chars", http.StatusBadRequest)
		return
	}

	usersMx.Lock()
	users[len(users)+1] = *user
	usersMx.Unlock()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User created successfully")
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	resBody, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, string(resBody))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	paramId := r.PathValue("id")
	userId, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Please provide integer userId in params", http.StatusBadRequest)
		return
	}

	usersMx.RLock()
	user, ok := users[userId]
	usersMx.RUnlock()
	if !ok {
		http.Error(w, fmt.Sprintf("User with id %d is not exist", userId), http.StatusBadRequest)
		return
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error in parsing user to json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(jsonUser))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	paramId := r.PathValue("id")
	userId, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "Please provide integer userId in params", http.StatusBadRequest)
		return
	}

	usersMx.RLock()
	_, ok := users[userId]
	usersMx.RUnlock()
	if !ok {
		http.Error(w, fmt.Sprintf("User with id %d is not exist", userId), http.StatusBadRequest)
		return
	}

	usersMx.Lock()
	delete(users, userId)
	usersMx.Unlock()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User deleted successfully")
}
