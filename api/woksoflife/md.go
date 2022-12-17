package woksoflife

import (
	"fmt"
	"net/http"
	"strings"

	"humblebundle-scraper/internal"
)

//goland:noinspection GoUnusedExportedFunction
func Handler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	if !queryParams.Has("url") {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("the query params 'url'"))
		return
	}

	url := queryParams.Get("url")

	recipe := internal.GetRecipe(url)

	response := strings.Builder{}
	response.WriteString(
		fmt.Sprintf(
			`# %s\n\n<img src="%s" width="700"/>\n> tempo de preparo: %v\n## ingredientes\n`,
			recipe.Name,
			recipe.Image,
			recipe.PrepTime,
		),
	)

	for _, ingredient := range recipe.Ingredients {
		response.WriteString(fmt.Sprintf("* %s\n", ingredient))
	}

	response.WriteString("\n## modo de preparo\n")
	for i, instruction := range recipe.Instructions {
		response.WriteString(fmt.Sprintf("%d. %s\n", i+1, instruction))
	}

	response.WriteString("Observação: " + recipe.Notes)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(response.String()))
}
