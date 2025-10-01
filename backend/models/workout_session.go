package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutSession struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID  `bson:"userId" json:"userId"`
	RoutineID *primitive.ObjectID `bson:"routineId,omitempty" json:"routineId,omitempty"` // optional
	StartTime time.Time           `bson:"startTime" json:"startTime"`
	EndTime   time.Time           `bson:"endTime" json:"endTime"`
	Notes     string              `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
}
