package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutineExercise struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoutineID   primitive.ObjectID `bson:"routineId" json:"routineId"`
	ExerciseID  primitive.ObjectID `bson:"exerciseId" json:"exerciseId"`
	Orden       int                `bson:"orden" json:"orden"`
	Series      int                `bson:"series" json:"series"`
	Repeticiones int               `bson:"repeticiones" json:"repeticiones"`
	PesoObjetivo *float64          `bson:"pesoObjetivo,omitempty" json:"pesoObjetivo,omitempty"`
	TiempoObjSeg *int              `bson:"tiempoObjetivoSeg,omitempty" json:"tiempoObjetivoSeg,omitempty"`
	Notas       string             `bson:"notas,omitempty" json:"notas,omitempty"`
}
