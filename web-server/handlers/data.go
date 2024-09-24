package handlers

import "sync"

type User struct {
	Name string `json="name"`
}

var (
	users   = make(map[int]User)
	usersMx sync.RWMutex
)
