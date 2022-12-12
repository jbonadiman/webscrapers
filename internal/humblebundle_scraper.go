package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func GetBundleData(browserlessToken string, url string) (string, []string) {
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

	var bundleName string
	var items []*cdp.Node
	var itemNames []string

	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),
		chromedp.AttributeValue(
			".bundle-logo",
			"alt",
			&bundleName,
			nil,
			chromedp.ByQuery,
		),
		chromedp.Nodes(".item-bundleName", &items, chromedp.ByQueryAll),
	); err != nil {
		log.Fatalf("Failed data from %s: %v", url, err)
	}

	for _, node := range items {
		itemNames = append(
			itemNames,
			"- "+node.Children[0].NodeValue,
		)
	}

	return bundleName, itemNames
}
