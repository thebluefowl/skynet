package controller

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/thebluefowl/skynet/geo"
	"github.com/thebluefowl/skynet/model"
)

type Controller struct {
	e *echo.Echo
	m *model.MetricStore
}

func NewOverviewController(e *echo.Echo, m *model.MetricStore) *Controller {
	return &Controller{
		e: e,
		m: m,
	}
}

func (c *Controller) Register() {
	c.e.POST("/metrics", c.Create)
	c.e.GET("/overview", c.GetOverview)
	c.e.GET("/trends", c.GetTrends)
	c.e.GET("/stats/web", c.GetBrowserStats)
	c.e.GET("/stats/mobile", c.GetMobileStats)
	c.e.GET("/stats/routes", c.GetRouteStats)
	c.e.GET("/stats/countries", c.GetCountryStats)
}

type RequestMetric struct {
	TimeToFirstByte        float64 `json:"time_to_first_byte"`
	FirstPaint             float64 `json:"first_paint"`
	FirstContentfulPaint   float64 `json:"first_contentful_paint"`
	FirstInputDelay        float64 `json:"first_input_delay"`
	LargestContentfulPaint float64 `json:"largest_contentful_paint"`
	TotalBlockingTime      float64 `json:"total_blocking_time"`
	CumulativeLayoutShift  float64 `json:"cumulative_layout_shift"`
}

type Metadata struct {
	BrowserName         string `json:"browser_name"`
	BrowserVersion      string `json:"browser_version"`
	OperatingSystem     string `json:"operating_system"`
	NetworkInformation  string `json:"network_information"`
	DeviceMemory        string `json:"device_memory"`
	HardwareConcurrency string `json:"hardware_concurrency"`
	ServiceWorkerStatus string `json:"service_worker_status"`
	Route               string `json:"route"`
	IsLowEndDevice      bool   `json:"is_low_end_device"`
	IsMobileDevice      bool   `json:"is_mobile_device"`
	IsLowEndExperience  bool   `json:"is_low_end_experience"`
	IPAddress           string `json:"ip_address"`
}

type CreateMetricRequest struct {
	Metrics   RequestMetric `json:"metrics"`
	Metadata  Metadata      `json:"metadata"`
	Timestamp *time.Time    `json:"timestamp"`
}

func (ctrlr *Controller) Create(c echo.Context) error {
	request := &CreateMetricRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}

	ip := c.RealIP()
	if request.Metadata.IPAddress != "" {
		ip = request.Metadata.IPAddress
	}

	fmt.Println(ip)
	country, lat, long, err := geo.GetCountry(ip)
	if err != nil {
		fmt.Println(err)
		return c.JSON(400, err)
	}

	metric := &model.Metric{
		TimeToFirstByte:        request.Metrics.TimeToFirstByte,
		FirstPaint:             request.Metrics.FirstPaint,
		FirstContentfulPaint:   request.Metrics.FirstContentfulPaint,
		FirstInputDelay:        request.Metrics.FirstInputDelay,
		LargestContentfulPaint: request.Metrics.LargestContentfulPaint,
		TotalBlockingTime:      request.Metrics.TotalBlockingTime,
		CumuLayoutShift:        request.Metrics.CumulativeLayoutShift,
		BrowserName:            request.Metadata.BrowserName,
		BrowserVersion:         request.Metadata.BrowserVersion,
		OperatingSystem:        request.Metadata.OperatingSystem,
		NetworkInformation:     request.Metadata.NetworkInformation,
		DeviceMemory:           request.Metadata.DeviceMemory,
		HardwareConcurrency:    request.Metadata.HardwareConcurrency,
		ServiceWorkerStatus:    request.Metadata.ServiceWorkerStatus,
		LocationLatitude:       lat,
		LocationLongitude:      long,
		Country:                country,
		Route:                  request.Metadata.Route,
		IsLowEndDevice:         request.Metadata.IsLowEndDevice,
		IsMobileDevice:         request.Metadata.IsMobileDevice,
		IsLowEndExperience:     request.Metadata.IsLowEndExperience,
		Timestamp:              request.Timestamp,
	}

	if err := ctrlr.m.Create(metric); err != nil {
		fmt.Println("error creating metric", err)
		return c.JSON(500, err)
	}
	return c.JSON(200, request)
}

func (ctrlr *Controller) GetOverview(c echo.Context) error {
	response, err := ctrlr.m.GetOverview("time_to_first_byte")
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(200, response)
}

func (ctrlr *Controller) GetTrends(c echo.Context) error {
	response, err := ctrlr.m.GetTrends("time_to_first_byte")
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(200, response)
}

func (ctrlr *Controller) GetBrowserStats(c echo.Context) error {
	response, err := ctrlr.m.GetWebBrowserStats()
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(200, response)
}

func (ctrlr *Controller) GetMobileStats(c echo.Context) error {
	response, err := ctrlr.m.GetMobileStats()
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(200, response)
}

func (ctrlr *Controller) GetRouteStats(c echo.Context) error {
	response, err := ctrlr.m.GetRouteStats()
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(200, response)
}

func (ctrlr *Controller) GetCountryStats(c echo.Context) error {
	response, err := ctrlr.m.GetCountryStats()
	if err != nil {
		fmt.Println(err)
	}
	return c.JSON(200, response)
}
