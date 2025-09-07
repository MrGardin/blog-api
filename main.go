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
	router.GET("/getUsers", handler.GetAllUsers)
	router.GET("/getUser/:id", handler.GetUserById)

	http.ListenAndServe(":8080", router)
}
