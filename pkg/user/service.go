package userpkg

import (
	"context"
	"errors"
	"log"
	"mahi-go-explorer/internal/config"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Service defines interface for the user service
type Service interface {
	EnsureAdminUserExists() error
	LoginUser(req *LoginRequest) (string, error)

	CreateUser(user *User) (any, error)
	GetUsers(conds bson.M, opts *options.FindOptions) ([]User, error)
	GetUser(conds bson.M, opts *options.FindOneOptions) (*User, error)
	UpdateUser(conds bson.M, update bson.M, opts *options.UpdateOptions) (any, error)
	DeleteUser(conds bson.M) (any, error)
}

type service struct {
	db   *mongo.Database
	coll *config.Collection
}

// NewService returns new instance of user service
func NewService(db *mongo.Database, coll *config.Collection) Service {
	return service{db, coll}
}

func (s service) EnsureAdminUserExists() error {
	//check if any user exists
	count, err := s.db.Collection(s.coll.UserCollection).CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	adminUser := User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@example.com",
		Role:      "ADMIN",
	}

	hp, _ := hashPassword("admin123")
	adminUser.HashedPassword = hp

	_, err = s.CreateUser(&adminUser)
	if err != nil {
		return err
	}

	log.Println("Default admin user created: admin@example.com / admin123")

	return nil
}

func (s service) LoginUser(req *LoginRequest) (string, error) {
	var user User
	err := s.db.Collection(s.coll.UserCollection).FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return "", errors.New("user not found")
	}

	if user.IsBlocked {
		return "", errors.New("user is blocked")
	}

	if !camparePassword(user.HashedPassword, req.Password) {
		return "", errors.New("invalid password")
	}

	//create a jwt
	claims := JwtClaims{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		Exp:       time.Now().Add(1 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//sign & return the token
	t, err := token.SignedString([]byte(config.GetFromEnv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (s service) CreateUser(user *User) (any, error) {
	resp, err := s.db.Collection(s.coll.UserCollection).InsertOne(context.TODO(), user, nil)
	if err != nil {
		return nil, err
	}
	return resp.InsertedID, nil
}

func (s service) GetUsers(conds bson.M, opts *options.FindOptions) ([]User, error) {
	var users []User
	cursor, err := s.db.Collection(s.coll.UserCollection).Find(context.TODO(), conds, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s service) GetUser(conds bson.M, opts *options.FindOneOptions) (*User, error) {
	var user User
	err := s.db.Collection(s.coll.UserCollection).FindOne(context.TODO(), conds, opts).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s service) UpdateUser(conds bson.M, update bson.M, opts *options.UpdateOptions) (any, error) {
	resp, err := s.db.Collection(s.coll.UserCollection).UpdateOne(context.TODO(), conds, update, opts)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s service) DeleteUser(conds bson.M) (any, error) {
	resp, err := s.db.Collection(s.coll.UserCollection).DeleteOne(context.TODO(), conds)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
