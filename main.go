package main

import (
	"blog-api/handlers"
	"blog-api/storage"
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	connString := "postgres://postgres:admin@localhost:5432/blog-api?sslmode=disable"

	//Создаем пул соединений
	pool, err := storage.NewPool(context.Background(), connString)
	if err != nil {
		log.Fatalf("Error with pool connection")
	}
	defer pool.Close()
	log.Println("Success")
	log.Println("http://localhost:8080")

	handler := handlers.NewHandler(pool)
	router := httprouter.New()

	router.GET("/", handlers.HomeHandler)
	router.GET("/version", handler.CheckVersionDataBase)
	router.GET("/users", handler.GetAllUsers)
	router.GET("/users/:id", handler.GetUserById)
	router.POST("/users", handler.CreateUser)

	http.ListenAndServe(":8080", router)
}
