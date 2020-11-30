package handlers

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
	"truth/database"
)

type AuthResponse struct {
	Token string `json:"token"`
}

type LoginBody struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Register godoc
// @Tags Auth
// @Summary Register
// @Accept  json
// @Produce json
// @Param user body database.User true "User"
// @Success 200 {object} AuthResponse
// @Router /auth/register [post]
func Register(c echo.Context) error {
	userBody := new(database.User)
	userBody.Active = true
	userBody.RoleID = 2
	userBody.CreatedAt = time.Now()
	userBody.UpdatedAt = time.Now()

	if err := c.Bind(userBody); err != nil {
		return err
	}
	if err := c.Validate(userBody); err != nil {
		return err
	}

	u, _ := database.FindUserByEmail(userBody.Email)
	if u != nil {
		return echo.ErrBadRequest
	}

	user, err := database.CreateUser(userBody)
	if err != nil {
		log.Println("error", err)
		return err
	}

	token, err := user.GenerateJWT()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, bson.M{
		"token": token,
	})
}

// Login godoc
// @Tags Auth
// @Summary Login
// @Accept  json
// @Produce json
// @Param user body LoginBody true "LoginBody"
// @Success 200 {object} AuthResponse
// @Router /auth/login [post]
func Login(c echo.Context) error {
	body := new(LoginBody)
	if err := c.Bind(body); err != nil {
		return err
	}
	if err := c.Validate(body); err != nil {
		return err
	}

	user, _ := database.FindUserByEmail(body.Email)
	if user == nil {
		return echo.ErrUnauthorized
	}

	ok := user.CheckPassword(body.Password)
	if !ok {
		return echo.ErrUnauthorized
	}

	token, err := user.GenerateJWT()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, bson.M{
		"token": token,
	})
}
