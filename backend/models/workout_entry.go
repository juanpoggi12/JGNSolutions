package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutEntry struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WorkoutSessionID primitive.ObjectID `bson:"workoutSessionId" json:"workoutSessionId"`
	ExerciseID       primitive.ObjectID `bson:"exerciseId" json:"exerciseId"`
	Serie            int                `bson:"serie" json:"serie"`
	RepsHechas       *int               `bson:"repsHechas,omitempty" json:"repsHechas,omitempty"`
	PesoUsado        *float64           `bson:"pesoUsado,omitempty" json:"pesoUsado,omitempty"`
	TiempoSeg        *int               `bson:"tiempoSeg,omitempty" json:"tiempoSeg,omitempty"`
	PercepcionEsfuerzo *int             `bson:"percepcionEsfuerzo,omitempty" json:"percepcionEsfuerzo,omitempty"` // RIR/PE opcional
}