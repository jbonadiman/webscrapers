package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo"
)

var browserlessApiKey string

func main() {
	browserlessApiKey = os.Getenv("BROWSERLESS_API_KEY")
	if browserlessApiKey == "" {
		panic("the BROWSERLESS_API_KEY environment variable is not set!")
	}

	e := echo.New()
	e.GET("/md", getBundleMarkdown)

	e.Logger.Fatal(e.Start(":8080"))
}

func getBundleMarkdown(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.String(
			http.StatusBadRequest,
			"the query param 'url' is required",
		)
	}

	allocatorContext, cancel := chromedp.NewRemoteAllocator(
		context.Background(),
		fmt.Sprintf("wss://chrome.browserless.io?token=%s", browserlessApiKey),
	)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	var title string
	var items []*cdp.Node
	var itemNames []string

	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.AttributeValue(
			".bundle-logo",
			"alt",
			&title,
			nil,
			chromedp.ByQuery,
		),
		chromedp.Nodes(".item-title", &items, chromedp.ByQueryAll),
	); err != nil {
		log.Fatalf("Failed getting title of example.com: %v", err)
	}

	for _, node := range items {
		itemNames = append(
			itemNames,
			"- "+node.Children[0].NodeValue,
		)
	}

	return c.String(
		http.StatusOK, fmt.Sprintf(
			"Humble Bundle \"%s\"\n\n%s",
			strings.TrimLeft(strings.Split(title, ":")[1], " "),
			strings.Join(itemNames, "\n"),
		),
	)
}
