package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           primitive.ObjectID `bson:"userId" json:"userId"`
	RefreshTokenHash string             `bson:"refreshTokenHash" json:"-"`
	ExpiresAt        time.Time          `bson:"expiresAt" json:"expiresAt"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	RevokedAt        *time.Time         `bson:"revokedAt,omitempty" json:"revokedAt,omitempty"`
}
