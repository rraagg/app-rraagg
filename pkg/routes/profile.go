package routes

import (
	"github.com/rraagg/rraagg/pkg/controller"
	"github.com/rraagg/rraagg/templates"

	"github.com/labstack/echo/v4"
)

type (
	profile struct {
		controller.Controller
	}
)

func (c *profile) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = templates.LayoutMain
	page.Name = templates.PageProfile
	page.Metatags.Description = "Welcome to my profile page."
	page.Metatags.Keywords = []string{"Go", "MVC", "Web", "Software"}
	page.Pager = controller.NewPager(ctx, 4)
	page.Data = "This is my profile page."

	return c.RenderPage(ctx, page)
}
