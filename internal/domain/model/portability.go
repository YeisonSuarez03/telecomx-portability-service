package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// PortabilityDetails includes user-specific portability data.
type PortabilityDetails struct {
	Number string `bson:"number" json:"number"`
	Notes  string `bson:"notes,omitempty" json:"notes"`
}

// Portability is the root document stored in MongoDB.
type Portability struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"user_id" json:"userID"`
	OperatorInfo  string             `bson:"operator_info" json:"operatorInfo"`
	CurrentStatus string             `bson:"current_status" json:"currentStatus"` // e.g. Requested, InProgress, Completed, Suspended, Failed
	Details       PortabilityDetails `bson:"details" json:"details"`
	RequestedAt   time.Time          `bson:"requested_at" json:"requestedAt"`
	CompletedAt   *time.Time         `bson:"completed_at" json:"completedAt"`
}
