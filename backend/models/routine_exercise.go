package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoutineExercise struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoutineID     primitive.ObjectID `bson:"routineId" json:"routineId"`
	ExerciseID    primitive.ObjectID `bson:"exerciseId" json:"exerciseId"`
	Order         int                `bson:"order" json:"order"`
	Sets          int                `bson:"sets" json:"sets"`
	Reps          int                `bson:"reps" json:"reps"`
	TargetWeight  *float64           `bson:"targetWeight,omitempty" json:"targetWeight,omitempty"`
	TargetTimeSec *int               `bson:"targetTimeSec,omitempty" json:"targetTimeSec,omitempty"`
	Notes         string             `bson:"notes,omitempty" json:"notes,omitempty"`
}
