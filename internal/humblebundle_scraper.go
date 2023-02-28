package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Bundle struct {
	Name  string
	Items []string
}

func getBundleContent(browserlessToken string, url string) ([]byte, error) {
	reqBody, err := json.Marshal(
		map[string]string{"url": url},
	)

	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		fmt.Sprintf(
			"https://chrome.browserless.io/content?token=%s&headless=true&blockAds=true",
			browserlessToken,
		),
		"application/json",
		bytes.NewBuffer(reqBody),
	)

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetBundleData(browserlessToken string, url string) (
	Bundle,
	error,
) {
	htmlContent, err := GrabContent(browserlessToken, url)
	if err != nil {
		return Bundle{}, nil
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		return Bundle{}, err
	}

	bundleName, _ := doc.Find(".bundle-logo").First().Attr("alt")

	var itemNames []string

	items := doc.Find(".item-title")
	for _, bundleItem := range items.Nodes {
		itemNames = append(
			itemNames,
			fmt.Sprintf("- %s", bundleItem.FirstChild.Data),
		)
	}

	return Bundle{
		Name:  bundleName,
		Items: itemNames,
	}, nil
}
