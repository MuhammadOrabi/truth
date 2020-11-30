package crawler

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
)

type CNNCard struct {
	Title string `json:"title"`
	URL string `json:"url"`
	Image string `json:"image"`
	Date string `json:"date"`
	Description string `json:"description"`
}

func (craw Crawler) StartCNN() ([]CNNCard, error) {
	url := craw.Source.URL + "?q=" + craw.Keyword

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return nil, fmt.Errorf("could not navigate: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.WaitVisible(`.cnn-search__result.cnn-search__result--article`)); err != nil {
		return nil, fmt.Errorf("could not navigate: %v", err)
	}

	err := chromedp.Run(ctx, chromedp.ActionFunc(func(context.Context) error {
		log.Printf(">>>>>>>>>>>>>>>>>>>> element IS VISIBLE")
		return nil
	}))
	if err != nil {
		return nil, fmt.Errorf("could not do action: %v", err)
	}

	var titles []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(".cnn-search__result-headline > a", &titles)); err != nil {
		return nil, fmt.Errorf("could not get titles: %v", err)
	}

	response := make([]CNNCard, len(titles))

	var images []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(".cnn-search__result-thumbnail > a > img", &images)); err != nil {
		return nil, fmt.Errorf("could not get images: %v", err)
	}

	var descriptions []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(".cnn-search__result-contents > .cnn-search__result-body", &descriptions)); err != nil {
		return nil, fmt.Errorf("could not get descriptions: %v", err)
	}


	var dates []*cdp.Node
	if err := chromedp.Run(ctx, chromedp.Nodes(".cnn-search__result-contents > .cnn-search__result-publish-date > span:nth-child(2)", &dates)); err != nil {
		return nil, fmt.Errorf("could not get dates: %v", err)
	}

	for i, d := range titles {
		img := images[i]
		desc := descriptions[i]
		date := dates[i]
		response[i].Image = img.Attributes[1]
		response[i].URL = d.Attributes[1]
		response[i].Title = d.Children[0].NodeValue
		response[i].Description = desc.Children[0].NodeValue
		response[i].Date = date.Children[0].NodeValue
	}

	return response, nil
}
