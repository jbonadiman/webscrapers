package function

import (
	"fmt"
	"net/http"
	"strings"

	"handler/function/sources"
)

const (
	Url              = "url"
	BrowserlessToken = "browserlessToken"
	Format           = "format"
)

func throwError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(
		[]byte(err.Error()),
	)
}

func throwIfInvalidQueryParams(w http.ResponseWriter, r *http.Request) error {
	queryParams := r.URL.Query()

	if !queryParams.Has(Url) ||
		!queryParams.Has(BrowserlessToken) ||
		!queryParams.Has(Format) {

		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(
			[]byte(fmt.Sprintf(
				"query params %q, %q and %q are required",
				Url,
				BrowserlessToken,
				Format,
			)),
		)

		return fmt.Errorf("invalid params")
	}

	return nil
}

func Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := throwIfInvalidQueryParams(w, r)
	if err != nil {
		return
	}

	queryParams := r.URL.Query()
	url := queryParams.Get(Url)

	if strings.HasPrefix(url, "https://www.humblebundle.com") {
		bundle, err := sources.GetBundleData(
			queryParams.Get(BrowserlessToken),
			queryParams.Get(Url),
		)
		if err != nil {
			throwError(w, err)
			return
		}

		switch queryParams.Get(Format) {
		case "json":
			jsonBundle, err := bundle.ToJSON()
			if err != nil {
				throwError(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonBundle)
			break

		case "md":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "text/markdown")
			_, _ = w.Write(bundle.ToMD())
			break

		default:
			throwError(
				w,
				fmt.Errorf("invalid format %q", queryParams.Get(Format)),
			)
			return
		}
	}

	if strings.HasPrefix(url, "https://www.woksoflife.com") {
		throwError(w, fmt.Errorf("not implemented yet"))
		return
	}

	w.Header().Set("Cache-Control", "max-age=0, s-maxage=86400")
}
