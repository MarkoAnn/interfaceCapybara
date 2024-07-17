package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// User представляет сущность пользователя
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	Create(user User) error
	List() ([]User, error)
	Find(id string) (User, error)
	Update(user User) error
	Delete(id string) error
}

// InMemoryUserRepository реализует интерфейс UserRepository с использованием in-memory хранилища
type InMemoryUserRepository struct {
	users map[string]User
	mu    sync.Mutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]User)}
}

func (repo *InMemoryUserRepository) Create(user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.users[user.ID]; exists {
		return errors.New("user already exists")
	}
	repo.users[user.ID] = user
	return nil
}

func (repo *InMemoryUserRepository) List() ([]User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	var userList []User
	for _, user := range repo.users {
		userList = append(userList, user)
	}
	return userList, nil
}

func (repo *InMemoryUserRepository) Find(id string) (User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	user, exists := repo.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

func (repo *InMemoryUserRepository) Update(user User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	repo.users[user.ID] = user
	return nil
}

func (repo *InMemoryUserRepository) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if _, exists := repo.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(repo.users, id)
	return nil
}

// HTTP Handlers

func createUserHandler(repo UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")
		ageStr := r.URL.Query().Get("age")

		if id == "" || name == "" || ageStr == "" {
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			return
		}

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Invalid age parameter", http.StatusBadRequest)
			return
		}

		user := User{ID: id, Name: name, Age: age}
		if err := repo.Create(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func listUsersHandler(repo UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := repo.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func findUserHandler(repo UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		user, err := repo.Find(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}

func updateUserHandler(repo UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")
		ageStr := r.URL.Query().Get("age")

		if id == "" || name == "" || ageStr == "" {
			http.Error(w, "Missing parameters", http.StatusBadRequest)
			return
		}

		age, err := strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Invalid age parameter", http.StatusBadRequest)
			return
		}

		user := User{ID: id, Name: name, Age: age}
		if err := repo.Update(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func deleteUserHandler(repo UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing id parameter", http.StatusBadRequest)
			return
		}

		if err := repo.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	// Создаем экземпляр InMemoryUserRepository
	userRepository := NewInMemoryUserRepository()

	http.HandleFunc("/create", createUserHandler(userRepository))
	http.HandleFunc("/list", listUsersHandler(userRepository))
	http.HandleFunc("/find", findUserHandler(userRepository))
	http.HandleFunc("/update", updateUserHandler(userRepository))
	http.HandleFunc("/delete", deleteUserHandler(userRepository))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
