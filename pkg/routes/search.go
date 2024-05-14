package routes

import (
	"fmt"
	"math/rand"

	"github.com/rraagg/rraagg/pkg/controller"
	"github.com/rraagg/rraagg/templates"

	"github.com/labstack/echo/v4"
)

type (
	search struct {
		controller.Controller
	}

	searchResult struct {
		Title string
		URL   string
	}
)

func (c *search) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = templates.LayoutMain
	page.Name = templates.PageSearch

	// Fake search results
	var results []searchResult
	if search := ctx.QueryParam("query"); search != "" {
		for i := 0; i < 5; i++ {
			title := "Lorem ipsum example ddolor sit amet"
			index := rand.Intn(len(title))
			title = title[:index] + search + title[index:]
			results = append(results, searchResult{
				Title: title,
				URL:   fmt.Sprintf("https://www.%s.com", search),
			})
		}
	}
	page.Data = results

	return c.RenderPage(ctx, page)
}
