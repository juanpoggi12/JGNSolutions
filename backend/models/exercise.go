package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExerciseCategory string
type GrupoMuscular string
type Dificultad string

const (
	CategoriaFuerza       ExerciseCategory = "FUERZA"
	CategoriaCardio       ExerciseCategory = "CARDIO"
	CategoriaFlexibilidad ExerciseCategory = "FLEXIBILIDAD"
	CategoriaOtra         ExerciseCategory = "OTRA"

	GrupoPecho    GrupoMuscular = "PECHO"
	GrupoEspalda  GrupoMuscular = "ESPALDA"
	GrupoPierna   GrupoMuscular = "PIERNA"
	GrupoHombro   GrupoMuscular = "HOMBRO"
	GrupoBrazo    GrupoMuscular = "BRAZO"
	GrupoCore     GrupoMuscular = "CORE"

	DifBaja  Dificultad = "BAJA"
	DifMedia Dificultad = "MEDIA"
	DifAlta  Dificultad = "ALTA"
)

type Exercise struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Nombre          string             `bson:"nombre" json:"nombre"`
	Descripcion     string             `bson:"descripcion,omitempty" json:"descripcion,omitempty"`
	Categoria       ExerciseCategory   `bson:"categoria" json:"categoria"`
	GrupoMuscular   GrupoMuscular      `bson:"grupoMuscular" json:"grupoMuscular"`
	Dificultad      Dificultad         `bson:"dificultad" json:"dificultad"`
	MediaURL        string             `bson:"mediaUrl,omitempty" json:"mediaUrl,omitempty"`
	Instrucciones   []string           `bson:"instrucciones,omitempty" json:"instrucciones,omitempty"`
	CreatedByUserID primitive.ObjectID `bson:"createdByUserId" json:"createdByUserId"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsDeleted       bool               `bson:"isDeleted" json:"isDeleted"`
}
