package main

import (
	"log"

	"humblebundle-scraper/internal"
)

func main() {
	recipe := internal.GetRecipe(
		"https://thewoksoflife.com/jasmine-shortbread-cookies/",
	)

	log.Println(recipe)
	log.Println(recipe.PrepTime)
}
