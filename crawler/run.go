package crawler

import "truth-finder/database"

type Crawler struct {
	Source database.Source
	Keyword string
}

func (craw *Crawler) Run() (interface{}, error) {
	if craw.Source.Name == "CNN" {
		return craw.StartCNN()
	}
	return nil, nil
}
