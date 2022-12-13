package main

import (
	"os"

	"humblebundle-scraper/internal"
)

func main() {
	_, _, _ = internal.GetBundleData(
		os.Getenv("BROWSERLESS_TOKEN"),
		"",
	)
}
