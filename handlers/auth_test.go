package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"truth/config"
	"truth/model"
	"truth/storage"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// CustomValidator ...
type CustomValidator struct {
	validator *validator.Validate
}

// Validate ...
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
func TestRegister(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	conString := config.GetTestingPostgresConnectionString()
	db := storage.NewDB(conString)
	db.AutoMigrate(&model.User{}, &model.Source{})

	db.Delete(&model.User{})
	userJSON, _ := json.Marshal(map[string]interface{}{
		"name":     "Muhammad Orabi",
		"email":    "orabi@mail.com",
		"password": "password",
		"role_id":  1,
	})
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Assertions
	if assert.NoError(t, Register(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestLogin(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	conString := config.GetTestingPostgresConnectionString()
	db := storage.NewDB(conString)
	db.AutoMigrate(&model.User{}, &model.Source{})
	userJSON, _ := json.Marshal(map[string]string{
		"email":    "orabi@mail.com",
		"password": "password",
	})
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	// Assertions
	if assert.NoError(t, Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
