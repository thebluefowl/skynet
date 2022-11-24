package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/thebluefowl/skynet/controller"
	"github.com/thebluefowl/skynet/model"
)

func SeedData() {

	// request := controller.CreateMetricRequest{}
}

type Curve struct {
	Height           float64
	Peak             float64
	StandarDeviation float64
}

func (c *Curve) GetVal(x float64, min float64, Height float64) float64 {
	g := Height * math.Exp(-math.Pow(x-c.Peak, 2)/(2*math.Pow(c.StandarDeviation, 2)))
	y := g + min
	r := rand.Float64()*2 - 1
	y = float64(y) * (1 + r)
	return y
}

func main() {

	browsers := model.MobileBrowsers

	c := Curve{
		Height:           2200,
		Peak:             5,
		StandarDeviation: 10,
	}

	startTime := time.Now().Add((-24 * 30) * time.Hour)
	endTime := time.Now()

	for startTime.Before(endTime) {
		//generate a random number between 1 and 60
		d := rand.Intn(60) + 1
		randomBrowser := browsers[rand.Intn(len(browsers))]
		startTime = startTime.Add(time.Duration(d) * time.Minute)
		request := controller.CreateMetricRequest{
			Metrics: controller.RequestMetric{
				TimeToFirstByte:        c.GetVal(float64(d), 200, 2200),
				FirstPaint:             c.GetVal(float64(d), 300, 1000),
				FirstContentfulPaint:   c.GetVal(float64(d), 1200, 1800),
				FirstInputDelay:        c.GetVal(float64(d), 20, 480),
				LargestContentfulPaint: c.GetVal(float64(d), 2000, 3500),
				TotalBlockingTime:      c.GetVal(float64(d), 100, 200),
				CumulativeLayoutShift:  c.GetVal(float64(d), 0.1, 0.5),
			},

			Metadata: controller.Metadata{
				BrowserName: randomBrowser,
				// random boolean
				IsMobileDevice: rand.Intn(2) == 1,
			},
			Timestamp: &startTime,
		}

		// make http post with json request
		data, err := json.Marshal(request)
		if err != nil {
			panic(err)
		}
		_, err = http.Post("http://localhost:1323/metrics", "application/json", bytes.NewBuffer(data))
		if err != nil {
		}
		fmt.Println(request.Metadata.BrowserName)

	}

}
