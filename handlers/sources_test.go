package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"truth/config"
	"truth/model"
	"truth/storage"

	"github.com/labstack/echo/v4/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestListSources(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	conString := config.GetTestingPostgresConnectionString()
	db := storage.NewDB(conString)
	db.AutoMigrate(&model.User{}, &model.Source{})
	userJSON, _ := json.Marshal(map[string]interface{}{
		"email":    "orabi@mail.com",
		"password": "password",
	})
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(userJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	var data map[string]string
	// Assertions
	if assert.NoError(t, Login(c)) {
		err = json.Unmarshal(rec.Body.Bytes(), &data)
		if err != nil {
			log.Println("error:", err)
		}
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	handler := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &model.JwtClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		BeforeFunc: func(context echo.Context) {
			token := context.Request().Header.Get("Authorization")
			if !strings.HasPrefix(token, "Bearer") {
				context.Request().Header.Set("Authorization", "Bearer "+token)
			}
		},
	})(ListSources)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, data["token"])
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateSource(t *testing.T) {
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
	var data map[string]string
	// Assertions
	if assert.NoError(t, Login(c)) {
		err = json.Unmarshal(rec.Body.Bytes(), &data)
		if err != nil {
			log.Println("error:", err)
		}
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	handler := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &model.JwtClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		BeforeFunc: func(context echo.Context) {
			token := context.Request().Header.Get("Authorization")
			if !strings.HasPrefix(token, "Bearer") {
				context.Request().Header.Set("Authorization", "Bearer "+token)
			}
		},
	})(CreateSource)
	sourceJSON, _ := json.Marshal(map[string]interface{}{
		"active": true,
		"name": "string",
		"status": "string",
		"tags": [1]string{"string"},
		"url": "string",
	})
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(sourceJSON)))
	req.Header.Set(echo.HeaderAuthorization, data["token"])
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	if assert.NoError(t, handler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
