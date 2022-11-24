package controller

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/thebluefowl/skynet/model"
)

type OverviewController struct {
	e *echo.Echo
	m *model.MetricStore
}

func NewOverviewController(e *echo.Echo, m *model.MetricStore) *OverviewController {
	return &OverviewController{
		e: e,
		m: m,
	}
}

func (c *OverviewController) Register() {
	c.e.GET("/overview", c.GetOverview)
}

func (ctrlr *OverviewController) GetOverview(c echo.Context) error {
	response, err := ctrlr.m.GetOverview("time_to_first_byte")
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(200, response)
}
