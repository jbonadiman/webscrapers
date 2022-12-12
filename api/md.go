package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func getBundleData(browserlessToken string, url string) (string, []string) {
	allocatorContext, cancel := chromedp.NewRemoteAllocator(
		context.Background(),
		fmt.Sprintf(
			"wss://chrome.browserless.io?token=%s",
			browserlessToken,
		),
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
		log.Fatalf("Failed data from %s: %v", url, err)
	}

	for _, node := range items {
		itemNames = append(
			itemNames,
			"- "+node.Children[0].NodeValue,
		)
	}

	return title, itemNames
}

//goland:noinspection GoUnusedExportedFunction
func Handler(w http.ResponseWriter, r *http.Request) {
	browserlessApiKey := os.Getenv("BROWSERLESS_API_KEY")
	if browserlessApiKey == "" {
		panic("the BROWSERLESS_API_KEY environment variable is not set!")
	}

	queryParams := r.URL.Query()

	if !queryParams.Has("url") {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("the query param 'url' is required"))
		return
	}
	url := queryParams.Get("url")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Cache-Control", "max-age=0, s-maxage=86400")

	bundleName, bundleItems := getBundleData(browserlessApiKey, url)

	response := fmt.Sprintf(
		"Humble Bundle \"%s\" (%d items)\n\n%s",
		strings.TrimLeft(strings.Split(bundleName, ":")[1], " "),
		len(bundleItems),
		strings.Join(bundleItems, "\n"),
	)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}
