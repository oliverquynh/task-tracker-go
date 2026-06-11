package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
)

type Task struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   int    `json:"created_at"`
	UpdatedAt   int    `json:"updated_at"`
}

func createJsonFile() error {
	filePath := "tasks.json"

	// O_CREATE: Create the file if it doesn't exist
	// O_EXCL: Used with O_CREATE, file must not exist
	// O_WRONLY: Open the file for writing only
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)

	if err != nil {
		// Check if the error is specifically because the file already exists
		if errors.Is(err, os.ErrExist) {
			// fmt.Println("File already exists. Doing nothing.")
			return nil
		}
	}

	err = os.WriteFile(filePath, []byte("[]"), 0666)

	defer file.Close()

	return err
}

func readTasks() ([]Task, error) {
	jsonStr, err := os.ReadFile("tasks.json")

	if err != nil {
		return nil, err
	}

	var tasks []Task

	err = json.Unmarshal(jsonStr, &tasks)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func writeTasks(tasks []Task) error {
	file, err := os.Create("tasks.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "    ")

	err = encoder.Encode(tasks)

	if err != nil {
		return err
	}

	return nil
}

func printTaskRow(task Task, idMaxLen int, descMaxLen int, statusMaxLen int) {
	idLen := len(strconv.Itoa(int(task.ID)))
	descLen := len(task.Description)
	statusLen := len(task.Status)

	idStr := fmt.Sprintf("%d", task.ID) + strings.Repeat(" ", idMaxLen - idLen)
	descStr := task.Description + strings.Repeat(" ", descMaxLen - descLen)
	statusStr := task.Status + strings.Repeat(" ", statusMaxLen - statusLen)

	t := time.Unix(int64(task.CreatedAt), 0)
	c := t.Format("2006-01-02 15:04:05")
	t = time.Unix(int64(task.UpdatedAt), 0)
	u := t.Format("2006-01-02 15:04:05")

	fmt.Printf("%s  %s  %s  %s  %s\n", idStr, descStr, statusStr, c, u)
}

func printTask(task Task) {
	t := time.Unix(int64(task.CreatedAt), 0)
	c := t.Format("2006-01-02 15:04:05")
	t = time.Unix(int64(task.UpdatedAt), 0)
	u := t.Format("2006-01-02 15:04:05")

	fmt.Printf("- ID: %d\n", task.ID)
	fmt.Printf("- Description: %s\n", task.Description)
	fmt.Printf("- Status: %s\n", task.Status)
	fmt.Printf("- Created at: %s\n", c)
	fmt.Printf("- Updated at: %s\n", u)
}

// task-tracker list
func listHandler(tasks []Task, argv []string) int {
	if len(tasks) == 0 {
		fmt.Println("[INFO] There are no tasks right now.")
		return 0
	}
	var idMaxLen int = len("ID")
	var descMaxLen int = len("Description")
	var statusMaxLen int = len(StatusInProgress)
	for _, task := range tasks {
		idLen := len(strconv.Itoa(int(task.ID)))
		descLen := len(task.Description)
		if idLen > idMaxLen {
			idMaxLen = idLen
		}
		if descLen > descMaxLen {
			descMaxLen = descLen
		}
	}

	idCol := "ID" + strings.Repeat(" ", idMaxLen - 2)
	descCol := "Description" + strings.Repeat(" ", descMaxLen - 11)
	statusCol := "Status" + strings.Repeat(" ", statusMaxLen - 6)
	createdAtCol := "Created at" + strings.Repeat(" ", 19 - 10) // 2026-06-11 00:00:00
	updatedAtCol := "Updated at" + strings.Repeat(" ", 19 - 10) // 2026-06-11 00:00:00

	fmt.Printf("%s  %s  %s  %s  %s\n", idCol, descCol, statusCol, createdAtCol, updatedAtCol)

	for _,task := range tasks {
		printTaskRow(task, idMaxLen, descMaxLen, statusMaxLen)
	}
	return 0
}

// task-tracker add "This is a sample task"
func addHandler(tasks []Task, argv []string) int {
	if len(argv) < 3 {
		fmt.Println("[ERROR] Please provide the description.")
		return 1
	}

	tasks, err := readTasks()

	if err != nil {
		fmt.Printf("[ERROR] Failed to read tasks: %v\n", err)
		return 1
	}

	var max uint = 0

	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}

	id := max + 1

	description := argv[2]

	task := Task{
		ID:          id,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   int(time.Now().Unix()),
		UpdatedAt:   int(time.Now().Unix()),
	}

	tasks = append(tasks, task)

	err = writeTasks(tasks)

	if err != nil {
		fmt.Printf("[ERROR] Failed to write tasks: %v\n", err)
		return 1
	}

	fmt.Println("[INFO] Task added successfully.")

	printTask(task)

	return 0
}

// task-tracker edit 1 "Change task description"
func editHandler(tasks []Task, argv []string) int {
	if len(argv) < 4 {
		fmt.Println("[ERROR] Please provide the task ID and description.")
		return 1
	}

	v, err := strconv.Atoi(argv[2])

	if err != nil || v <= 0 {
		fmt.Println("[ERROR] Task ID is invalid. It's must be an positive integer.")
		return 1
	}

	id := uint(v)

	var found bool
	var task Task

	for i := range tasks {
		if tasks[i].ID == id {
			description := argv[3]
			tasks[i].Description = description
			tasks[i].UpdatedAt = int(time.Now().Unix())
			found = true
			task = tasks[i]
		}
	}

	if !found {
		fmt.Printf("[ERROR] Task ID [%d] was not found.\n", id)
		return 1
	}

	err = writeTasks(tasks)

	if err != nil {
		fmt.Printf("[ERROR] Failed to write tasks: %v\n", err)
		return 1
	}

	fmt.Printf("[INFO] Task [%d] edited successfully.\n", id)

	printTask(task)

	return 0
}

// task-tracker mark 1 in-progress
func markHandler(tasks []Task, argv []string) int {
	if len(argv) < 4 {
		fmt.Println("[ERROR] Please provide the task ID and status.")
		return 1
	}

	v, err := strconv.Atoi(argv[2])

	if err != nil || v <= 0 {
		fmt.Println("[ERROR] Task ID is invalid. It's must be an positive integer.")
		return 1
	}

	id := uint(v)

	status := argv[3] // TODO: Validate if status is valid.

	if status != StatusTodo && status != StatusInProgress && status != StatusDone {
		fmt.Printf("[ERROR] Task status [%s] is invalid. Please choose one of %s, %s or %s.\n", status, StatusTodo, StatusInProgress, StatusDone)
		return 1
	}

	var found bool
	var task Task

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = int(time.Now().Unix())
			task = tasks[i]
			found = true
		}
	}

	if !found {
		fmt.Printf("[ERROR] Task ID [%d] was not found.\n", id)
		return 1
	}

	err = writeTasks(tasks)

	if err != nil {
		fmt.Printf("[ERROR] Failed to write tasks: %v\n", err)
		return 1
	}

	fmt.Printf("[INFO] Task [%d] marked as [%s] successfully.\n", id, status)

	printTask(task)

	return 0
}

// task-tracker delete 1
func deleteHandler(tasks []Task, argv []string) int {
	if len(argv) < 3 {
		fmt.Println("[ERROR] Please provide the task ID.")
		return 1
	}

	v, err := strconv.Atoi(argv[2])

	if err != nil || v <= 0 {
		fmt.Println("[ERROR] Task ID is invalid. It's must be an positive integer.")
		return 1
	}

	id := uint(v)

	var found bool

	var copied []Task

	for i := range tasks {
		if tasks[i].ID == id {
			found = true
			continue
		}

		copied = append(copied, tasks[i])
	}

	if !found {
		fmt.Printf("[ERROR] Task ID [%d] was not found.\n", id)
		return 1
	}

	err = writeTasks(copied)

	if err != nil {
		fmt.Printf("[ERROR] Failed to write tasks: %v\n", err)
		return 1
	}

	fmt.Printf("[INFO] Task [%d] deleted successfully.\n", id)


	return 0
}

func main() {
	err := createJsonFile()

	if err != nil {
		fmt.Println("[ERROR] Could not create the tasks.json file")
	}

	if len(os.Args) == 1 {
		fmt.Println("[ERROR] Please provide a command.")
		os.Exit(1)
	}

	command := os.Args[1]

	tasks, err := readTasks()

	if err != nil {
		fmt.Printf("[ERROR] Failed to read tasks: %v\n", err)
		os.Exit(1)
	}

	var handler func(tasks []Task, argv []string) int

	switch command {
	case "list":
		handler = listHandler
	case "add":
		handler = addHandler
	case "edit":
		handler = editHandler
	case "mark":
		handler = markHandler
	case "delete":
		handler = deleteHandler
	default:
		fmt.Printf("[ERROR] command %s was not found.\n", command)
		os.Exit(1)
	}

	exitCode := handler(tasks, os.Args)
	os.Exit(exitCode)
}
