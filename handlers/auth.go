package handlers

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
	"truth/model"
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
// @Param user body model.User true "User"
// @Success 200 {object} AuthResponse
// @Router /auth/register [post]
func Register(c echo.Context) error {
	userBody := new(model.User)
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

	u := model.FindUserByEmail(userBody.Email)
	if u.ID != 0 {
		return echo.ErrBadRequest
	}

	user := model.CreateUser(userBody)

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

	user := model.FindUserByEmail(body.Email)
	if user.ID == 0 {
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
