package main

import (
	"context"
	"fmt"
	"github.com/LayssonENS/go-genealogy-api/config"
	"github.com/LayssonENS/go-genealogy-api/database"
	_ "github.com/LayssonENS/go-genealogy-api/docs"
	personHttpDelivery "github.com/LayssonENS/go-genealogy-api/person/delivery/http"
	personRepository "github.com/LayssonENS/go-genealogy-api/person/repository"
	personUCase "github.com/LayssonENS/go-genealogy-api/person/usecase"
	relationshipHttpDelivery "github.com/LayssonENS/go-genealogy-api/relationships/delivery/http"
	relationshipsRepository "github.com/LayssonENS/go-genealogy-api/relationships/repository"
	relationshipUCase "github.com/LayssonENS/go-genealogy-api/relationships/usecase"
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

// @title Go Genealogy API
// @version 1.0
// @description This is Genealogy API in Go.

// @host 0.0.0.0:8000
// @BasePath /auth
// @schemes http
// @query.collection.format multi

func main() {
	ctx := context.Background()
	log := logrus.New()

	dbInstance, err := database.NewPostgresConnection(config.GetEnv().DbConfig)
	if err != nil {
		log.WithError(err).Fatal("failed connection database")
		return
	}

	err = database.DBMigrate(dbInstance, config.GetEnv().DbConfig)
	if err != nil {
		log.WithError(err).Fatal("failed to migrate")
		return
	}

	router := gin.Default()

	pRepository := personRepository.NewPostgresPersonRepository(dbInstance)
	personService := personUCase.NewPersonUseCase(pRepository)

	rRepository := relationshipsRepository.NewPostgresRelationshipRepository(dbInstance)
	relationshipService := relationshipUCase.NewRelationshipUseCase(rRepository)

	personHttpDelivery.NewPersonHandler(router, personService)
	relationshipHttpDelivery.NewRelationshipHandler(router, relationshipService)
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
