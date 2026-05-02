// @title           Subscription Service API
// @version         1.0
// @description     Subscription Management
// @host            localhost:8080
// @BasePath        /
package main

import (
	_ "github.com/trippleflipp/test-task-go/docs"
	"github.com/trippleflipp/test-task-go/internal/db"
	"github.com/trippleflipp/test-task-go/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Info("Starting logger")

	database := db.InitDB()
	defer database.Close()

	h := handlers.NewSubscriptionHandler(database, log)

	r := gin.Default()

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.POST("/subscriptions", h.Create)
		api.PUT("/subscriptions/:id", h.Update)
		api.GET("/subscriptions", h.List)
		api.GET("/subscriptions/:id", h.GetByID)
		api.DELETE("/subscriptions/:id", h.Delete)

		api.GET("/subscriptions/total", h.GetTotal)
	}

	log.Info("Server started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error running server: %v", err)
	}
}
