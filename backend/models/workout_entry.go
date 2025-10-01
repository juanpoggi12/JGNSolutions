package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WorkoutEntry struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WorkoutSessionID primitive.ObjectID `bson:"workoutSessionId" json:"workoutSessionId"`
	ExerciseID       primitive.ObjectID `bson:"exerciseId" json:"exerciseId"`
	SetNumber        int                `bson:"setNumber" json:"setNumber"`
	RepsDone         *int               `bson:"repsDone,omitempty" json:"repsDone,omitempty"`
	WeightUsed       *float64           `bson:"weightUsed,omitempty" json:"weightUsed,omitempty"`
	TimeSec          *int               `bson:"timeSec,omitempty" json:"timeSec,omitempty"`
	PerceivedEffort  *int               `bson:"perceivedEffort,omitempty" json:"perceivedEffort,omitempty"` // optional RPE
}
