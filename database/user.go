package database

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"time"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" swaggerignore:"true"`
	Name      string             `json:"name" bson:"name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	Active    bool               `json:"active" bson:"active" validate:"required" swaggerignore:"true"`
	RoleID    int                `json:"role_id" bson:"role_id" validate:"required" swaggerignore:"true"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" swaggerignore:"true"`
}

type JwtClaims struct {
	User User `json:"user"`
	jwt.StandardClaims
}

func hashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func (u *User) CheckPassword(plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(u.Password)
	bytePlainPwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (u *User) GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = u
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// CreateUser ...
func CreateUser(user *User) (*User, error) {
	db := Load()

	user.Password = hashAndSalt([]byte(user.Password))

	res, err := db.Users.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("err", err)
		return nil, err
	}

	ID, _ := res.InsertedID.(primitive.ObjectID)
	user.ID = ID

	return user, nil
}

// FindUser ...
func FindUserByEmail(email string) (*User, error) {
	db := Load()

	var user User
	filter := bson.D{primitive.E{Key: "email", Value: email}}

	err := db.Users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
