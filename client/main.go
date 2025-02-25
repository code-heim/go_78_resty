package main

import (
	"fmt"
	"log"
	"time"

	"resty.dev/v3"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	client := resty.New().SetBaseURL("http://localhost:8080")

	// 1. Create a new task (POST)
	newTask := Task{Title: "Learn Resty", Done: false}
	var createdTask Task

	_, err := client.R().SetTimeout(2*time.Second).
		SetHeader("Content-Type", "application/json").
		SetBody(newTask).
		SetResult(&createdTask).
		Post("/tasks")

	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	fmt.Printf("Created Task: %+v\n", createdTask)

	// 2. Get all tasks (GET)
	var tasks []Task
	_, err = client.R().
		SetResult(&tasks).
		Get("/tasks")

	if err != nil {
		log.Fatalf("Failed to get tasks: %v", err)
	}

	fmt.Println("\nAll Tasks:")
	for _, task := range tasks {
		fmt.Printf("- ID: %d, Title: %s, Done: %t\n",
			task.ID, task.Title, task.Done)
	}

	// 3. Update a task (PUT)
	updatedTask := Task{Title: "Master Resty", Done: true}
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(updatedTask).
		SetResult(&updatedTask).
		Put(fmt.Sprintf("/tasks/%d", createdTask.ID))

	if err != nil {
		log.Fatalf("Failed to update task: %v", err)
	}

	fmt.Printf("\nUpdated Task: %+v\n", updatedTask)

	// 4. Delete a task (DELETE)
	_, err = client.R().
		Delete(fmt.Sprintf("/tasks/%d", createdTask.ID))

	if err != nil {
		log.Fatalf("Failed to delete task: %v", err)
	}

	fmt.Printf("\nDeleted Task with ID %d\n", createdTask.ID)
}
