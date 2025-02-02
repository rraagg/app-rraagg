package routes

import (
	"github.com/rraagg/rraagg/pkg/context"
	"github.com/rraagg/rraagg/pkg/controller"
	"github.com/rraagg/rraagg/pkg/services"
	"github.com/rraagg/rraagg/templates"

	"github.com/labstack/echo/v4"
)

type (
	weather struct {
		controller.Controller
		X      int
		Y      int
		Office string
	}

	geocodingLocationForm struct {
		Street     string `form:"street" validate:"required"`
		City       string `form:"city"   validate:"required"`
		State      string `form:"state"  validate:"required"`
		Submission controller.FormSubmission
	}
)

func (w *weather) Get(ctx echo.Context) error {
	page := controller.NewPage(ctx)
	page.Layout = templates.LayoutMain
	page.Name = templates.PageWeather
	page.Metatags.Description = "Welcome to the weather page."
	page.Metatags.Keywords = []string{"Go", "MVC", "Web", "Software"}
	page.Form = geocodingLocationForm{}
	page.Pager = controller.NewPager(ctx, 4)

	if form := ctx.Get(context.FormKey); form != nil {
		page.Form = form.(*geocodingLocationForm)
	}

	page.Data = w.getHourlyForecast(ctx, &page.Pager)
	return w.RenderPage(ctx, page)
}

func (w *weather) Post(ctx echo.Context) error {
	var form geocodingLocationForm
	ctx.Set(context.FormKey, &form)

	// Parse the form values
	if err := ctx.Bind(&form); err != nil {
		ctx.Logger().Error(err)
		return w.Fail(err, "unable to bind form")
	}
	if err := form.Submission.Process(ctx, form); err != nil {
		return w.Fail(err, "unable to process form submission")
	}

	coordinates, err := w.Container.Geocoding.GetGeocodingCoordinates(
		ctx,
		form.Street,
		form.City,
		form.State,
	)
	if err != nil {
		ctx.Logger().Error(err)
	}

	properties, err := w.Container.Points.GetPoints(ctx, coordinates.Y, coordinates.X)
	if err != nil {
		return w.Fail(err, "unable to get points")
	}

	w.X = properties.GridX
	w.Y = properties.GridY
	w.Office = properties.GridID

	ctx.Logger().Infof("Coordinates: X = %v, Y = %v", coordinates.X, coordinates.Y)

	return w.Get(ctx)
}

func (w *weather) getHourlyForecast(ctx echo.Context, pager *controller.Pager) []services.Period {
	ctx.Logger().Info("Fetching hourly forecast")
	if w.X == 0 {
		w.X = 74
	}
	if w.Y == 0 {
		w.Y = 89
	}
	if w.Office == "" {
		w.Office = "FGZ"
	}

	periods, err := w.Container.Weather.GetHourlyForecast(ctx, w.X, w.Y, w.Office)
	if err != nil {
		ctx.Logger().Error(err)
		errPeriods := make([]services.Period, 0)
		return errPeriods
	}
	pager.SetItems(len(periods))

	return periods[pager.GetOffset() : pager.GetOffset()+pager.ItemsPerPage]
}
