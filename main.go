package main

import (
	"fmt"
	"log"
	"os"
    "net/http"
    "github.com/a-h/templ"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
	"encoding/json"

	"github.com/Nader-Rahhal/tasky/handlers"
)

var svc *dynamodb.DynamoDB
var tableName string

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(os.Getenv("AWS_REGION")),
    })
    if err != nil {
        log.Fatalf("Failed to create session: %v", err)
    }
    svc = dynamodb.New(sess)
    tableName = os.Getenv("DYNAMODB_TABLE_NAME")
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/delete-task", deleteTaskHandler)
	http.HandleFunc("/add-task", addTaskHandler)
	http.HandleFunc("/tasks", getTasksHandler)

    fmt.Println("Listening on :3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    component := home()
    templ.Handler(component).ServeHTTP(w, r)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := r.URL.Query().Get("id")
	taskTitle := r.URL.Query().Get("title")

	if taskID == "" {
        http.Error(w, "Missing task ID", http.StatusBadRequest)
        return
    }

	if taskTitle == "" {
        http.Error(w, "Missing task title", http.StatusBadRequest)
        return
    }

	err := handlers.PutItem(svc, tableName, taskID, taskTitle)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding task: %v", err), http.StatusInternalServerError)
        return
	}

	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Task added successfully"})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    taskID := r.URL.Query().Get("id")
    if taskID == "" {
        http.Error(w, "Missing task ID", http.StatusBadRequest)
        return
    }

    err := handlers.DeleteTask(svc, tableName, taskID)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error deleting task: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	tasks, err := handlers.GetAllTableItems(svc, tableName)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks: %v", err), http.StatusInternalServerError)
        return
	}

	json.NewEncoder(w).Encode(tasks)


}