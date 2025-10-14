package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty"`
	UserID       primitive.ObjectID  `bson:"userId"`
	RefreshHash  string              `bson:"refreshHash"`
	ExpiresAt    time.Time           `bson:"expiresAt"`
	CreatedAt    time.Time           `bson:"createdAt"`
	RevokedAt    *time.Time          `bson:"revokedAt,omitempty"`
	UserAgent    string              `bson:"userAgent,omitempty"`
	IP           string              `bson:"ip,omitempty"`
	ReplacedByID *primitive.ObjectID `bson:"replacedById,omitempty"`
}
