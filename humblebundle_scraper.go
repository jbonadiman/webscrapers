package webscrapers

import (
	"bytes"
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
)

type Bundle struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}

func (b Bundle) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

func GetBundleData(browserlessToken string, url string) (
	Bundle,
	error,
) {
	htmlContent, err := GrabContent(browserlessToken, url)
	if err != nil {
		return Bundle{}, err
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
			bundleItem.FirstChild.Data,
		)
	}

	return Bundle{
		Name:  bundleName,
		Items: itemNames,
	}, nil
}
