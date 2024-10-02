package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/joho/godotenv"
)

type Task struct {
	ID    string `json:"ID"`
	Title string `json:"Title"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	svc := dynamodb.New(sess)
	fmt.Println("Successfully started session")

	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	err = putItem(svc, tableName, "2", "Complete Project")
	if err != nil {
		log.Printf("Error putting item: %v", err)
	}

	_, err = getAllTableItems(svc, tableName)
	if err != nil {
		log.Printf("Error getting all items: %v", err)
	}

	err = deleteTask(svc, tableName, "2")
	if err != nil {
		log.Printf("Error deleting task: %v", err)
	}

	_, err = getAllTableItems(svc, tableName)
	if err != nil {
		log.Printf("Error getting all items after deletion: %v", err)
	}
}

func putItem(svc *dynamodb.DynamoDB, tableName, id, title string) error {
	item := Task{
		ID:    id,
		Title: title,
	}
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("Got error marshalling map: %s", err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	return err
}

func getAllTableItems(svc *dynamodb.DynamoDB, tableName string) ([]Task, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	var tasks []Task
	var err error
	for {
		var result *dynamodb.ScanOutput
		result, err = svc.Scan(input)
		if err != nil {
			return nil, fmt.Errorf("got error scanning table: %s", err)
		}
		var items []Task
		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
		if err != nil {
			return nil, fmt.Errorf("got error unmarshalling items: %s", err)
		}
		tasks = append(tasks, items...)
		if result.LastEvaluatedKey == nil {
			break
		}
		input.ExclusiveStartKey = result.LastEvaluatedKey
	}
	fmt.Println(tasks)
	return tasks, nil
}

func deleteTask(svc *dynamodb.DynamoDB, tableName, key string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(key),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := svc.DeleteItem(input)
	if err != nil {
		return fmt.Errorf("failed to delete task %s: %w", key, err)
	}
	return nil
}