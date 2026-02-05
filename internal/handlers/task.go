package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks []Task
var nextID = 1


func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}


func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTasks(w, r)
	case http.MethodPost:
		createTask(w, r)
	case http.MethodPatch:
		updateTask(w, r)
	case http.MethodDelete: 
		deleteTask(w, r)
	default:
		errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			errorResponse(w, http.StatusBadRequest, "invalid id")
			return
		}
		for _, task := range tasks {
			if task.ID == id {
				jsonResponse(w, http.StatusOK, task)
				return
			}
		}
		errorResponse(w, http.StatusNotFound, "task not found")
		return
	}

	if tasks == nil {
		jsonResponse(w, http.StatusOK, []Task{})
		return
	}
	jsonResponse(w, http.StatusOK, tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil || newTask.Title == "" {
		errorResponse(w, http.StatusBadRequest, "invalid title")
		return
	}

	newTask.ID = nextID
	nextID++
	newTask.Done = false

	tasks = append(tasks, newTask)
	jsonResponse(w, http.StatusCreated, newTask)
}


func updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	var updateData struct {
		Done *bool `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil || updateData.Done == nil {
		errorResponse(w, http.StatusBadRequest, "done is required")
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = *updateData.Done
			jsonResponse(w, http.StatusOK, map[string]bool{"updated": true})
			return
		}
	}
	errorResponse(w, http.StatusNotFound, "task not found")
}


func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			jsonResponse(w, http.StatusOK, map[string]bool{"deleted": true})
			return
		}
	}

	errorResponse(w, http.StatusNotFound, "task not found")
}