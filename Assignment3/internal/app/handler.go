package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"assignment3/pkg/modules"
)

type UserHandler struct {
	usecase *UserUsecase
}

func NewUserHandler(usecase *UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(r.URL.Path, "/users")

	if path == "" || path == "/" {
		switch r.Method {
		case http.MethodGet:
			users, err := h.usecase.GetAllUsers()
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(users)
		case http.MethodPost:
			var user modules.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
				return
			}
			id, err := h.usecase.CreateUser(user)
			if err != nil {
				http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": id})
		default:
			http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		}
		return
	}

	idStr := strings.TrimPrefix(path, "/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID format"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.usecase.GetUserByID(id)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(user)
	case http.MethodPut:
		var user modules.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, `{"error": "Invalid input"}`, http.StatusBadRequest)
			return
		}
		err = h.usecase.UpdateUser(id, user)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"message": "User updated"})
	case http.MethodDelete:
		rows, err := h.usecase.DeleteUser(id)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "User deleted", "rows_affected": rows})
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}
