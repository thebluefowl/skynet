package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/thebluefowl/skynet/controller"
	"github.com/thebluefowl/skynet/db"
	"github.com/thebluefowl/skynet/model"
)

func main() {
	_, err := db.GetDB()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.CORS())

	store := model.NewMetricStore()
	overviewController := controller.NewOverviewController(
		e,
		store,
	)

	overviewController.Register()
	e.Logger.Fatal(e.Start(":1323"))

}
