package crawler

import (
	"truth/model"
)

type Crawler struct {
	Source model.Source
	Keyword string
}

func (craw *Crawler) Run() ([]CNNCard, error) {
	if craw.Source.Name == "CNN" {
		return craw.StartCNN()
	}
	return nil, nil
}
