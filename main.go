package main

import (
	"blog-api/handler"
	"blog-api/repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Пробуем загрузить .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env файл не найден, используем переменные системы")
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//Формируем строку подключения
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	//Подключаемся к бд
	dbBlogApi, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer dbBlogApi.Close()

	userRepo := repository.NewUserRepository(dbBlogApi)

	if err = userRepo.CheckConnection(); err != nil {
		log.Fatal("Ошибка проверки коннекта")
	}

	_, err = userRepo.GetWelcomeMessage()
	if err != nil {
		log.Fatalf("Ошибка получения приветственного сообщения", err)
	}

	http.HandleFunc("/", handler.HomeHandler(dbBlogApi))
	fmt.Println("🚀 Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
