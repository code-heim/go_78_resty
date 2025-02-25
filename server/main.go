package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	tasks     = []Task{}
	idCounter = 1
	mu        sync.Mutex
)

func main() {
	http.HandleFunc("/tasks", handleTasks)
	http.HandleFunc("/tasks/", handleTaskByID)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var newTask Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "Invalid request payload",
				http.StatusBadRequest)
			return
		}
		mu.Lock()
		newTask.ID = idCounter
		idCounter++
		tasks = append(tasks, newTask)
		mu.Unlock()
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}

func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var updatedTask Task
		if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
			http.Error(w, "Invalid request payload",
				http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		for i, task := range tasks {
			if task.ID == id {
				tasks[i].Title = updatedTask.Title
				tasks[i].Done = updatedTask.Done
				json.NewEncoder(w).Encode(tasks[i])
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)

	case http.MethodDelete:
		mu.Lock()
		defer mu.Unlock()
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "Task deleted"}`))
				return
			}
		}
		http.Error(w, "Task not found", http.StatusNotFound)

	default:
		http.Error(w, "Method not allowed",
			http.StatusMethodNotAllowed)
	}
}
