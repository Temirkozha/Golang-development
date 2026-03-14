package handler

import (
	"assignment5/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	nameFilter := r.URL.Query().Get("name")
	orderBy := r.URL.Query().Get("order_by")

	response, err := h.repo.GetPaginatedUsers(page, pageSize, nameFilter, orderBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	u1, _ := strconv.Atoi(r.URL.Query().Get("user1"))
	u2, _ := strconv.Atoi(r.URL.Query().Get("user2"))

	friends, err := h.repo.GetCommonFriends(u1, u2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

func (h *UserHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int `json:"user_id"`
		FriendID int `json:"friend_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	h.repo.AddFriend(req.UserID, req.FriendID)
	w.WriteHeader(http.StatusCreated)
}
