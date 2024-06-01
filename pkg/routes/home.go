package routes

import (
	"fmt"

	"github.com/rraagg/rraagg/pkg/controller"
	"github.com/rraagg/rraagg/pkg/services"
	"github.com/rraagg/rraagg/templates"

	"github.com/labstack/echo/v4"
)

type (
	home struct {
		controller.Controller
	}

	post struct {
		Title string
		Body  string
	}
)

func (c *home) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = templates.LayoutMain
	page.Name = templates.PageHome
	page.Metatags.Description = "Welcome to the homepage."
	page.Metatags.Keywords = []string{"Go", "MVC", "Web", "Software"}
	page.Pager = controller.NewPager(ctx, 4)
	page.Data = c.getHourlyForecast(ctx, &page.Pager)

	return c.RenderPage(ctx, page)
}

// fetchPosts is an mock example of fetching posts to illustrate how paging works
func (c *home) fetchPosts(pager *controller.Pager) []post {
	pager.SetItems(20)
	posts := make([]post, 20)

	for k := range posts {
		posts[k] = post{
			Title: fmt.Sprintf("Post example #%d", k+1),
			Body:  fmt.Sprintf("Lorem ipsum example #%d ddolor sit amet, consectetur adipiscing elit. Nam elementum vulputate tristique.", k+1),
		}
	}
	return posts[pager.GetOffset() : pager.GetOffset()+pager.ItemsPerPage]
}

func (c *home) getHourlyForecast(ctx echo.Context, pager *controller.Pager) []services.Period {
	ctx.Logger().Info("Fetching hourly forecast")
	tmpX := 74
	tmpY := 89
	office := "FGZ"

	periods, err := c.Container.Weather.GetHourlyForecast(ctx, tmpX, tmpY, office)
	if err != nil {
		ctx.Logger().Error(err)
		errPeriods := make([]services.Period, 0)
		return errPeriods
	}
	pager.SetItems(len(periods))

	return periods[pager.GetOffset() : pager.GetOffset()+pager.ItemsPerPage]
}
