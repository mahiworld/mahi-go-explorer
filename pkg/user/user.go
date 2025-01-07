package userpkg

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User defines user schema
type User struct {
	ID             primitive.ObjectID `json:"id,ompiempty" bson:"_id,omitempty"`
	FirstName      string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName       string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Email          string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone          string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Role           string             `json:"role,omitempty" bson:"role,omitempty"`
	HashedPassword string             `json:"-" bson:"hashedPassword,omitempty"`
	IsBlocked      bool               `json:"isBlocked" bson:"isBlocked"`
}

// CreateRequest defines user create request
type CreateRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role,omitempty"`
	Password  string `json:"password,omitempty"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func camparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// CreateUser creates user
func (cu *CreateRequest) CreateUser() (*User, error) {
	u := &User{
		FirstName: cu.FirstName,
		LastName:  cu.LastName,
		Email:     cu.Email,
		Phone:     cu.Phone,
		Role:      cu.Role,
		IsBlocked: false,
	}

	hp, err := hashPassword(cu.Password)
	if err != nil {
		return nil, err
	}

	u.HashedPassword = hp

	return u, nil
}

// LoginRequest defines login request schema
type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// UpdateRequest defines user update request
type UpdateRequest struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role,omitempty"`
}
