package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jbonadiman/webscrapers/api/humblebundle/json"
	"github.com/jbonadiman/webscrapers/api/humblebundle/md"
)

func main() {
	r := gin.Default()
	r.GET(
		"/api/humblebundle/md", func(c *gin.Context) {
			md.Handler(
				c.Writer,
				c.Request,
			)
		},
	)

	r.GET(
		"/api/humblebundle/json", func(c *gin.Context) {
			json.Handler(
				c.Writer,
				c.Request,
			)
		},
	)
	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}
