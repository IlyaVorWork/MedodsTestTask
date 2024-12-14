package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"MedodsTestTask/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение .env
	err := godotenv.Load("deploy/.env")
	if err != nil {
		log.Fatalf("cannot open .env file: %sv", err)
	}

	// Подключение бд
	connStr := fmt.Sprintf("postgresql://%s:%s@db/%s?sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	pgDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot open postgres connection ", err)
	}
	defer pgDB.Close()

	// Инициализация Gin
	router := gin.New()
	router.Use(gin.Recovery())

	// Первый слой логики - слой БД
	provider := auth.NewProvider(pgDB)
	// Второй слой логики - слой бизнес логики
	service := auth.NewUserService(provider)
	// Третий слой логики - слой HTTP хендлеров
	handler := auth.NewHandler(service)

	// Маршруты
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/login", handler.Login)
		authRouter.POST("/refresh", handler.Refresh)
	}

	// Запуск приложения
	addr := ":" + os.Getenv("SERVICE_PORT")
	err = router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
