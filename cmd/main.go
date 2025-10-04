package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"

	"telecomx-portability-service/internal/application/service"
	"telecomx-portability-service/internal/config"
	"telecomx-portability-service/internal/infrastructure/adapter/event/listener"
	"telecomx-portability-service/internal/infrastructure/adapter/repository"
	"telecomx-portability-service/internal/infrastructure/adapter/rest"
)

func main() {
	cfg := config.InstanceConfig()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("telecomx_portability")

	repo := repository.NewMongoRepository(db)
	svc := service.NewPortabilityService(repo)

	go func() {
		for {
			err := startKafkaWithRetry(svc, cfg)
			if err != nil {
				log.Println("Kafka listener failed, retrying in 5s:", err)
				time.Sleep(5 * time.Second)
				continue
			}
			break
		}
	}()

	mux := http.NewServeMux()
	rest.NewPortabilityHandler(svc).RegisterRoutes(mux)

	fmt.Printf("REST server running on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}

func startKafkaWithRetry(svc *service.PortabilityService, cfg config.Config) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in Kafka listener:", r)
		}
	}()

	return listener.StartKafkaListener(svc, cfg.Brokers, cfg.Topic, cfg.Group, cfg.Client)
}
