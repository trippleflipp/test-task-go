package main

import (
	"github.com/trippleflipp/test-task-go/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Info("Starting logger")

	database := db.InitDB()
	defer database.Close()

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		log.Info("GET/ping")
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	log.Info("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
