package listener

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"strconv"
	"time"

	"telecomx-portability-service/internal/application/service"
	"telecomx-portability-service/internal/domain/model"
)

type CustomerEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Phone struct {
	CodeNumber  string `json:"codeNumber,omitempty"`
	PhoneNumber int    `json:"phoneNumber,omitempty"`
}

type UserPayload struct {
	UserID    string `json:"userId"`
	Email     string `json:"email,omitempty"`
	Phone     Phone  `json:"phone"`
	Suspended bool   `json:"suspended,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`
}

func StartKafkaListener(svc *service.PortabilityService, brokers []string, topic, group, client string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: group,
		Dialer: &kafka.Dialer{
			ClientID: client,
		},
	})
	defer reader.Close()

	log.Printf("[Kafka] Listening on topic: %s", topic)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka error:", err)
			continue
		}

		var event CustomerEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Invalid event:", err)
			continue
		}

		var payload UserPayload
		_ = json.Unmarshal(event.Payload, &payload)

		switch event.Type {
		case "Customer.Created":
			err := svc.Create(context.Background(), &model.Portability{
				UserID:        payload.UserID,
				OperatorInfo:  "Telecomx",
				CurrentStatus: "InProgress",
				Details: model.PortabilityDetails{
					Number: strconv.Itoa(payload.Phone.PhoneNumber),
					Notes:  "Portability request from Telecomx",
				},
				RequestedAt: time.Now(),
			})
			if err != nil {
				log.Println("Error creating customer:", err)
				return
			}
		case "Customer.Updated":
			if payload.Deleted || payload.Suspended {
				err := svc.UpdateStatus(context.Background(), payload.UserID, "Cancelled")
				if err != nil {
					log.Println("Error updating customer status:", err)
					return
				}
			}
		case "Customer.Deleted":
			err := svc.Delete(context.Background(), payload.UserID)
			if err != nil {
				log.Println("Error deleting customer:", err)
				return
			}
		}
	}
}
