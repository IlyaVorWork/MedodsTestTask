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
	err := godotenv.Load("deploy/.env")
	if err != nil {
		log.Fatalf("cannot open .env file: %sv", err)
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@db/%s?sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	pgDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot open postgres connection ", err)
	}
	defer pgDB.Close()

	router := gin.New()
	router.Use(gin.Recovery())

	provider := auth.NewProvider(pgDB)
	service := auth.NewUserService(provider)
	handler := auth.NewHandler(service)
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/login", handler.Login)
		authRouter.POST("/refresh", handler.Refresh)
	}

	addr := ":" + os.Getenv("SERVICE_PORT")
	err = router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
