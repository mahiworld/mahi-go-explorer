package userpkg

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JwtClaims jwt claims
type JwtClaims struct {
	ID        primitive.ObjectID `json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Role      string             `json:"role"`
	Email     string             `json:"email"`
	Exp       interface{}        `json:"exp,omitempty"`
	jwt.StandardClaims
}

// UserContext user context
type UserContext struct {
	ID        primitive.ObjectID `json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Role      string             `json:"role"`
	Email     string             `json:"email"`
	Exp       interface{}        `json:"exp,omitempty"`
}
