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

	routes := []string{"/homepage",
		"/security",
		"/schedule-demo/",
		"/platform/",
		"/platform/autofix/",
		"/platform/code-coverage/",
		"/platform/iac-analysis/",
		"/platform/reports/",
		"/platform/sast/",
		"/platform/self-hosted/",
		"/platform/static-analysis/",
		"/changelog/",
		"/changelog/new-reports/",
		"/changelog/page/1/",
		"/changelog/sso-and-lcov/"}
	ips := []string{"116.240.147.252",
		"148.18.75.226",
		"30.222.237.37",
		"238.1.126.6",
		"109.251.255.118",
		"128.26.58.74",
		"51.42.159.123",
		"27.111.28.89",
		"171.185.137.228",
		"11.230.23.197",
		"119.22.105.211",
		"171.45.40.187",
		"126.65.23.17",
		"34.192.229.70",
		"251.207.226.152",
		"218.194.105.18",
		"217.119.90.189",
		"227.254.96.230",
		"237.77.191.212",
		"163.2.150.56",
		"61.11.254.119",
		"102.55.134.222",
		"3.154.88.170",
		"140.21.213.114",
		"98.122.203.180",
		"64.39.192.174",
		"60.72.163.137",
		"200.219.35.5",
		"20.248.97.164",
		"219.158.94.140",
		"132.183.76.163",
		"146.107.246.228",
		"203.192.247.242",
		"203.192.247.242",
		"203.192.247.242",
		"203.192.247.242",
		"203.192.247.242",
		"203.192.247.242",
		"203.192.247.242",
		"203.192.247.242",
		"43.211.170.221",
		"33.249.43.1",
		"141.230.206.97",
		"246.197.59.234",
		"46.187.115.31",
		"165.23.49.160",
		"182.102.62.251",
		"25.94.65.132",
		"122.95.104.19",
		"153.65.185.62",
		"68.132.85.154",
		"28.28.182.194",
		"42.251.115.179",
		"203.192.247.242"}

	for startTime.Before(endTime) {
		//generate a random number between 1 and 60
		d := rand.Intn(60) + 1
		randomBrowser := browsers[rand.Intn(len(browsers))]
		randomRoute := routes[rand.Intn(len(routes))]
		randomIP := ips[rand.Intn(len(ips))]
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
				BrowserName:    randomBrowser,
				Route:          randomRoute,
				IPAddress:      randomIP,
				IsMobileDevice: rand.Intn(2) == 1,
			},
			Timestamp: &startTime,
		}

		// make http post with json request
		data, err := json.Marshal(request)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))

		resp, err := http.Post("https://api.skynet.hackday.live/metrics", "application/json", bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.StatusCode)

	}

}
