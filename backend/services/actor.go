package services

import "go.mongodb.org/mongo-driver/bson/primitive"

// Actor representa al usuario autenticado que realiza una acción en el sistema.
type Actor struct {
	UserID primitive.ObjectID
	Role   string // "admin" o "user"
}
