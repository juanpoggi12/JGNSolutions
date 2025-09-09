package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.orgmongo-driver/bson/primitive"
)

type Nivel string
type Objetivo string

const (
	NivelPrincipiante Nivel = "PRINCIPIANTE"
	NivelIntermedio   Nivel = "INTERMEDIO"
	NivelAvanzado     Nivel = "AVANZADO"

	ObjetivoPerderPeso   Objetivo = "PERDER_PESO"
	ObjetivoGanarMusculo Objetivo = "GANAR_MUSCULO"
	ObjetivoMantenerse   Objetivo = "MANTENERSE"
)

type UserProfile struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `bson:"userId" json:"userId"`
	Nombre          string             `bson:"nombre" json:"nombre"`
	FechaNacimiento time.Time          `bson:"fechaNacimiento" json:"fechaNacimiento"`
	PesoKg          float64            `bson:"pesoKg" json:"pesoKg"`
	AlturaCm        int                `bson:"alturaCm" json:"alturaCm"`
	Nivel           Nivel              `bson:"nivel" json:"nivel"`
	Objetivo        Objetivo           `bson:"objetivo" json:"objetivo"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
