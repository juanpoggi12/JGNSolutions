package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"` // unique
	Email        string             `bson:"email" json:"email"`       // unique
	PasswordHash string             `bson:"passwordHash" json:"-"`    // never exposed in JSON
	Role         Role               `bson:"role" json:"role"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`
	IsActive     bool               `bson:"isActive" json:"isActive"`
}
