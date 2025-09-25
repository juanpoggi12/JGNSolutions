package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutSession struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	RoutineID      *primitive.ObjectID `bson:"routineId,omitempty" json:"routineId,omitempty"` // opcional
	FechaHoraInicio time.Time         `bson:"fechaHoraInicio" json:"fechaHoraInicio"`
	FechaHoraFin    time.Time         `bson:"fechaHoraFin" json:"fechaHoraFin"`
	NotasGenerales  string            `bson:"notasGenerales,omitempty" json:"notasGenerales,omitempty"`
	CreatedAt       time.Time         `bson:"createdAt" json:"createdAt"`
}