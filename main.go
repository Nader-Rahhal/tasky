package main

import (
    "fmt"
    "log"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Task struct {
    ID    string `json:"id"`
    Title string `json:"title"`
}

func main() {
    awsRegion := os.Getenv("AWS_REGION")
    if awsRegion == "" {
        log.Fatal("AWS_REGION environment variable is not set")
    }

    // Initialize a session
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(awsRegion),
    })
    if err != nil {
        log.Fatalf("Failed to create session: %v", err)
    }

    // Create DynamoDB client
    svc := dynamodb.New(sess)

    tableName := "Tasks"
    err = putItem(svc, tableName, "1", "Complete Project")
    if err != nil {
        log.Fatalf("Got error calling PutItem: %s", err)
    }
    fmt.Println("Successfully put item")
}

func putItem(svc *dynamodb.DynamoDB, tableName, id, title string) error {
    item := Task{
        ID:    id,
        Title: title,
    }

    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        return fmt.Errorf("got error marshalling map: %s", err)
    }

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(tableName),
    }

    _, err = svc.PutItem(input)
    return err
}