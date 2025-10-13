package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	FullName  string             `bson:"fullName" json:"fullName"`
	BirthDate time.Time          `bson:"birthDate" json:"birthDate"`
	WeightKg  float64            `bson:"weightKg" json:"weightKg"`
	HeightCm  int                `bson:"heightCm" json:"heightCm"`
	Level     Nivel              `bson:"level" json:"level"`
	Goal      Objetivo           `bson:"goal" json:"goal"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
