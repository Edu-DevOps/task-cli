package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"encoding/json"
	"os"
	"io/ioutil"

	"task-cli/models"
)


var (
    addAction bool
	actionDescription string
	status string
	id int
    updateAction bool
    deleteAction bool
	listAction bool
	filePath string = "tasks.json"
)

func main() {
	flag.BoolVar(&addAction, "add", false, "Add a task")
	flag.BoolVar(&updateAction, "update", false, "Update a task")
	flag.BoolVar(&deleteAction, "delete", false, "Delete a task")
	flag.BoolVar(&listAction, "list", false, "List tasks")
	// Used for add and update actions
	flag.StringVar(&actionDescription, "description", "My Task", "Description of the task")
	// Used for update and list actions
	flag.StringVar(&status, "status", "", "Status of the task")
	// Used for update and delete actions
	flag.IntVar(&id, "id", 0, "Id of the task")
	flag.Parse()

	_, err := os.Stat(filePath)
	if err != nil {
		createJson(filePath)
		fmt.Println("CREAT")
	}

	fmt.Println(actionSelector(addAction, updateAction, deleteAction))
}

func createJson (path string) {
	_, err := os.Create(path)
	if err != nil {
		log.Fatalln("Could not create the json file.")
	}
}

func actionSelector (addAction bool, updateAction bool, deleteAction bool) (string, error) {
	if addAction {
		req, err := add(actionDescription, filePath)
		if err != nil {
			return "", fmt.Errorf("Error adding the task: %v", err)
		}
		return req, nil
	} else if updateAction {
		req, err := update(id, actionDescription, status, filePath)
		if err != nil {
			return "", fmt.Errorf("Error updating the task: %v", err)
		}
		return req, nil
	} else if deleteAction {
		req, err := delete(id, filePath)
		if err != nil {
			return "", fmt.Errorf("Error deleting the task: %v", err)
		}
		return req, nil
	} else if listAction {
		req, err := list(filePath, status)
		if err != nil {
			return "", fmt.Errorf("Error listing the task: %v", err)
		}
		return req, nil
	} else {
		return "", fmt.Errorf("No action provided.")
	}
	return "", nil
}

// readJSONFile reads a JSON file and unmarshals it into a Task struct
func readJSONFile(filename string) ([]models.Task, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []models.Task{}, nil // Return an empty slice if the file doesn't exist
	}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read file content
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// If file is empty, return an empty slice
	if len(bytes) == 0 {
		return []models.Task{}, nil
	}

	// Unmarshal JSON into slice of tasks
	var tasks []models.Task
	if err := json.Unmarshal(bytes, &tasks); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return tasks, nil
}

func generateNextID(filename string) (int, error) {
	id := 1

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return id, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read file content
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return id, fmt.Errorf("error reading file: %v", err)
	}

	// If file is empty, return an empty slice
	if len(bytes) == 0 {
		return id, nil
	}

	// Unmarshal JSON into slice of tasks
	var tasks []models.Task
	if err := json.Unmarshal(bytes, &tasks); err != nil {
		return id, fmt.Errorf("error unmarshaling JSON: %v", err)
	}
	id = len(tasks)+1

	return id, nil
}

func add (description string, filePath string) (string, error) {
	currentTasks, err := readJSONFile(filePath)
	if err != nil {
		log.Fatalf("Could not read tasks: %v", err)
	}

	nextID, err := generateNextID(filePath)
	if err != nil {
		return "", fmt.Errorf("Could not generate next ID: %v", err)
	}
	// Create a new task
	newTask := models.Task{
		ID:          nextID,
		Description: description,
		Status:      "To-Do",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Append the new task to the slice
	currentTasks = append(currentTasks, newTask)
	jsonData, err := json.MarshalIndent(currentTasks, "", "\t")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}
	// Write data back to file
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing file: %v", err)
	}

	return "Task added successfully!", nil
}

func update (id int, description string, status string, filePath string) (string, error) {
	tasks, err := readJSONFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Could not read tasks: %v", err)
	}
	var updatedIndex int
	for index, task := range tasks {
		if task.ID == id {
			if description != "" {
				updatedIndex = index
				tasks[index].Description = description
			}
			if status != "" {
				updatedIndex = index
				tasks[index].Status = status
			}
			jsonData, err := json.MarshalIndent(tasks, "", "\t")
			if err != nil {
				return "", fmt.Errorf("error marshaling JSON: %v", err)
			}
			// Write data back to file
			err = ioutil.WriteFile(filePath, jsonData, 0644)
			if err != nil {
				return "", fmt.Errorf("error writing file: %v", err)
			}
		}
	}
	returnMessage := tasks[updatedIndex].Status + tasks[updatedIndex].Description
	return returnMessage, nil
}

func delete (id int, filePath string) (string, error) {
	tasks, err := readJSONFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Could not read tasks: %v", err)
	}
	index := id - 1
	tasks = append(tasks[:index], tasks[index+1:]...)
	jsonData, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %v", err)
	}
	// Write data back to file
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing file: %v", err)
	}
	returnMessage := fmt.Sprintf("Task %d deleted.", id)
	return returnMessage, nil
}

func list (filePath string, status string) (string, error) {
	tasks, err := readJSONFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Could not read tasks: %v", err)
	}
	var returnMessage string
	if status != "" {
		for _, task := range tasks {
			if status == task.Status {
				returnMessage = returnMessage + "\n" + task.Description + " " + task.Status
			}
		}
	} else {
		for _, task := range tasks {
			returnMessage = returnMessage + "\n" + task.Description + " " + task.Status
		}
	}

	return returnMessage, nil
}