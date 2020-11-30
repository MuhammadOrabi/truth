package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"truth-finder/database"

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
	claims := user.Claims.(*database.JwtClaims)
	log.Println("claims", *claims)
	sources, err := database.GetSources()
	if err != nil {
		return err
	}
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
	source := new(database.Source)
	if err := c.Bind(source); err != nil {
		return err
	}
	if err := c.Validate(source); err != nil {
		return err
	}
	if err := database.CreateSources(source); err != nil {
		return err
	}
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
	ID := c.Param("id")
	err := database.DeleteSource(ID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
