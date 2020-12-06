package handlers

import (
	"net/http"
	"truth/crawler"
	"truth/model"

	"github.com/labstack/echo/v4"
)

type searchBody struct {
	Keyword string `json:"keyword"`
	Type    string `json:"type"`
}

// Search godoc
// @Tags Search
// @Summary Search
// @Accept  json
// @Produce json
// @Param search body searchBody true "search params"
// @Success 200
// @Security ApiKeyAuth
// @Router /search [post]
func Search(c echo.Context) error {
	body := new(searchBody)
	if err := c.Bind(body); err != nil {
		return err
	}

	sources := model.GetSources()

	var data []interface{}
	for _, source := range sources {
		if source.Active {
			crawl := &crawler.Crawler{
				Source:  source,
				Keyword: body.Keyword,
			}

			output, err := crawl.Run()
			if err != nil {
				return err
			}

			data = append(data, output)
		}
	}
	return c.JSON(http.StatusOK, data)
}
