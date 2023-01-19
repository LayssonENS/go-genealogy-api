package main

import (
	"context"
	"fmt"
	personHttpDelivery "github.com/LayssonENS/go-genealogy-api/person/delivery/http"
	"github.com/LayssonENS/go-genealogy-api/person/repository"
	"github.com/LayssonENS/go-genealogy-api/person/usecase"
	"github.com/LayssonENS/go-genealogy-api/pkg/config"
	"github.com/LayssonENS/go-genealogy-api/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()

	log := logrus.New()
	//log.SetFormatter(&logrus.JSONFormatter{})

	dbInstance, err := database.NewPostgresConnection(config.GetEnv().DbConfig)
	if err != nil {
		log.WithError(err).Fatal("failed connection database")
		return
	}

	err = repository.DBMigrate(dbInstance, config.GetEnv().DbConfig)
	if err != nil {
		log.WithError(err).Fatal("failed to migrate")
		return
	}

	router := gin.Default()

	personRepository := repository.NewPostgresPersonRepository(dbInstance)
	userService := usecase.NewPersonUseCase(personRepository)

	personHttpDelivery.NewPersonHandler(router, userService)
	router.GET("/auth/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	gin.SetMode(gin.ReleaseMode)
	if config.GetEnv().Debug {
		gin.SetMode(gin.DebugMode)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.GetEnv().Port),
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down API...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("API Server forced to shutdown:", err)
	}

	log.Println("API Server exiting")
}
