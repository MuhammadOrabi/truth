package model

import (
	"log"
	"os"
	"time"
	"truth/storage"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint      `json:"id,omitempty" gorm:"primaryKey,autoIncrement" swaggerignore:"true"`
	Name      string    `json:"name" bson:"name" validate:"required"`
	Email     string    `json:"email" bson:"email" validate:"required"`
	Password  string    `json:"password" bson:"password" validate:"required"`
	Active    bool      `json:"active" bson:"active" validate:"required" swaggerignore:"true"`
	RoleID    int       `json:"role_id" bson:"role_id" validate:"required" swaggerignore:"true"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" swaggerignore:"true"`
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
func CreateUser(user *User) *User {
	db := storage.GetDBInstance()

	user.Password = hashAndSalt([]byte(user.Password))

	result := db.Create(&user)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return user
}

// FindUser ...
func FindUserByEmail(email string) *User {
	db := storage.GetDBInstance()
	var user User
	db.Where(&User{Email: email}).First(&user)

	return &user
}
