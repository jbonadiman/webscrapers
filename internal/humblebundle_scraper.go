package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

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
	string,
	[]string,
	error,
) {
	htmlContent, err := getBundleContent(browserlessToken, url)
	if err != nil {
		return "", nil, nil
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlContent))
	if err != nil {
		return "", nil, err
	}

	bundleName, _ := doc.Find(".bundle-logo").First().Attr("val")

	var itemNames []string

	items := doc.Find(".item-title")
	for _, bundleItem := range items.Nodes {
		itemNames = append(
			itemNames,
			fmt.Sprintf("- %s", bundleItem.FirstChild.Data),
		)
	}

	return bundleName, itemNames, nil
}
