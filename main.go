package main

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

type Task struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

var tasks = []Task{
    {ID: "1", Title: "clean sink", Completed: false},
    {ID: "2", Title: "clean room", Completed: false},
    {ID: "3", Title: "clean floor", Completed: false},
}

func getTasks(c echo.Context) error {
    return c.JSON(http.StatusOK, tasks)
}

func main() {
    e := echo.New()
    e.GET("/tasks", getTasks)
    e.Logger.Fatal(e.Start(":8080"))
}