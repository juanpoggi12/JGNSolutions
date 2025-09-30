package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Routine struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	Nombre      string             `bson:"nombre" json:"nombre"`
	Descripcion string             `bson:"descripcion,omitempty" json:"descripcion,omitempty"`
	IsTemplate  bool               `bson:"isTemplate,omitempty" json:"isTemplate,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsDeleted   bool               `bson:"isDeleted" json:"isDeleted"`
}

