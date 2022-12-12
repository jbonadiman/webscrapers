package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	browserlessApiKey := os.Getenv("BROWSERLESS_API_KEY")
	if browserlessApiKey == "" {
		panic("the BROWSERLESS_API_KEY environment variable is not set!")
	}

	url := flag.String(
		"url",
		"",
		"the Humble Book Bundle url",
	)

	devtoolsWsURL := flag.String(
		"devtools-ws-url",
		fmt.Sprintf("wss://chrome.browserless.io?token=%s", browserlessApiKey),
		"DevTools Websocket URL",
	)
	flag.Parse()

	if *url == "" {
		panic("the Humble Book Bundle url must be provided")
	}

	allocatorContext, cancel := chromedp.NewRemoteAllocator(
		context.Background(),
		*devtoolsWsURL,
	)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorContext)
	defer cancel()

	var title string
	var items []*cdp.Node
	var itemNames []string

	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(*url),
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

	fmt.Printf(
		"Humble Bundle \"%s\"\n\n%s",
		strings.TrimLeft(strings.Split(title, ":")[1], " "),
		strings.Join(itemNames, "\n"),
	)
}
