package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseCategory string
type MuscleGroup string
type Difficulty string

const (
	CategoriaFuerza       ExerciseCategory = "FUERZA"
	CategoriaCardio       ExerciseCategory = "CARDIO"
	CategoriaFlexibilidad ExerciseCategory = "FLEXIBILIDAD"
	CategoriaOtra         ExerciseCategory = "OTRA"

	MuscleChest    MuscleGroup = "PECHO"
	MuscleBack     MuscleGroup = "ESPALDA"
	MuscleLegs     MuscleGroup = "PIERNA"
	MuscleShoulder MuscleGroup = "HOMBRO"
	MuscleArm      MuscleGroup = "BRAZO"
	MuscleCore     MuscleGroup = "CORE"

	DiffLow  Difficulty = "BAJA"
	DiffMed  Difficulty = "MEDIA"
	DiffHigh Difficulty = "ALTA"
)

type Exercise struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`
	Description     string             `bson:"description,omitempty" json:"description,omitempty"`
	Category        ExerciseCategory   `bson:"category" json:"category"`
	MuscleGroup     MuscleGroup        `bson:"muscleGroup" json:"muscleGroup"`
	Difficulty      Difficulty         `bson:"difficulty" json:"difficulty"`
	MediaURL        string             `bson:"mediaUrl,omitempty" json:"mediaUrl,omitempty"`
	Instructions    []string           `bson:"instructions,omitempty" json:"instructions,omitempty"`
	CreatedByUserID primitive.ObjectID `bson:"createdByUserId" json:"createdByUserId"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsDeleted       bool               `bson:"isDeleted" json:"isDeleted"`
}
