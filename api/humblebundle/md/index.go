package md

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jbonadiman/webscrapers"
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

	bundle, err := webscrapers.GetBundleData(
		browserlessToken,
		url,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	mdResponse := strings.Builder{}
	mdResponse.WriteString(
		fmt.Sprintf(
			"Humble Bundle %q (%d items)\n\n",
			bundle.Name,
			len(bundle.Items),
		),
	)
	for _, item := range bundle.Items {
		mdResponse.WriteString(
			fmt.Sprintf(
				"- %s\n",
				item,
			),
		)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Add("Cache-Control", "max-age=0, s-maxage=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(mdResponse.String()))
}
