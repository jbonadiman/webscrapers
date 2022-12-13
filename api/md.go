package handler

import (
	"fmt"
	"net/http"
	"strings"

	"humblebundle-scraper/internal"
)

//goland:noinspection GoUnusedExportedFunction
func Handler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	if !queryParams.Has("url") || !queryParams.Has("browserlessToken") {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("the query params 'url' and 'browserlessToken' are required"))
		return
	}

	url := queryParams.Get("url")
	browserlessToken := queryParams.Get("browserlessToken")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Cache-Control", "max-age=0, s-maxage=86400")

	bundleName, bundleItems, err := internal.GetBundleData(
		browserlessToken,
		url,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	response := fmt.Sprintf(
		"Humble Bundle \"%s\" (%d items)\n\n%s",
		bundleName,
		len(bundleItems),
		strings.Join(bundleItems, "\n"),
	)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response))
}
