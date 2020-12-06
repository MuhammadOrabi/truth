package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strconv"
	"truth/model"

	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2/bson"
)

// ListSources godoc
// @Tags Sources
// @Summary List sources
// @Produce json
// @Success 200 {array} database.Source
// @Security ApiKeyAuth
// @Router /sources [get]
func ListSources(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*model.JwtClaims)
	log.Println("claims", *claims)
	sources := model.GetSources()
	return c.JSON(http.StatusOK, sources)
}

// CreateSource godoc
// @Tags Sources
// @Summary Create source
// @Accept  json
// @Produce json
// @Param source body database.Source true "Source"
// @Success 201
// @Security ApiKeyAuth
// @Router /sources [post]
func CreateSource(c echo.Context) error {
	source := new(model.Source)
	if err := c.Bind(source); err != nil {
		return err
	}
	if err := c.Validate(source); err != nil {
		return err
	}
	model.CreateSources(source)
	return c.JSON(http.StatusCreated, bson.M{})
}

// DeleteSource godoc
// @Tags Sources
// @Summary Delete source
// @Param id path string true "Source ID"
// @Success 200
// @Security ApiKeyAuth
// @Router /sources/{id} [delete]
func DeleteSource(c echo.Context) error {
	ID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	model.DeleteSource(uint(ID))
	return c.NoContent(http.StatusNoContent)
}
