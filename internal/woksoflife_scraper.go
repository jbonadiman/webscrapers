package internal

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Recipe struct {
	Name         string
	Image        string
	PrepTime     time.Duration
	Ingredients  []string
	Instructions []string
	Notes        string
}

func (r Recipe) String() string {
	recipe, _ := json.Marshal(r)
	return string(recipe)
}

func GetRecipe(url string) Recipe {
	recipe := Recipe{}

	c := colly.NewCollector()
	c.OnHTML(
		".wprm-recipe-name", func(e *colly.HTMLElement) {
			recipe.Name = e.Text
		},
	)

	c.OnHTML(
		"noscript > .attachment-post-thumbnail", func(e *colly.HTMLElement) {
			recipe.Image = e.Attr("src")
		},
	)

	c.OnHTML(
		".wprm-recipe-total_time-hours", func(e *colly.HTMLElement) {
			hours, err := strconv.Atoi(e.Text)
			if err == nil {
				recipe.PrepTime += time.Duration(hours) * time.Hour
			}
		},
	)

	c.OnHTML(
		".wprm-recipe-total_time-minutes", func(e *colly.HTMLElement) {
			minutes, err := strconv.Atoi(e.Text)
			if err == nil {
				recipe.PrepTime += time.Duration(minutes) * time.Minute
			}
		},
	)

	c.OnHTML(
		".wprm-recipe-ingredients", func(wrapper *colly.HTMLElement) {
			wrapper.ForEach(
				".wprm-recipe-ingredient",
				func(_ int, e *colly.HTMLElement) {
					recipe.Ingredients = append(
						recipe.Ingredients,
						strings.TrimLeft(e.Text, "▢ "),
					)
				},
			)
		},
	)

	c.OnHTML(
		".wprm-recipe-instructions", func(wrapper *colly.HTMLElement) {
			wrapper.ForEach(
				".wprm-recipe-instruction",
				func(_ int, e *colly.HTMLElement) {
					recipe.Instructions = append(
						recipe.Instructions,
						e.Text,
					)
				},
			)
		},
	)

	c.OnHTML(
		".wprm-recipe-notes", func(e *colly.HTMLElement) {
			recipe.Notes = e.Text
		},
	)

	_ = c.Visit(url)

	return recipe
}
