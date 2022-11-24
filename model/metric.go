package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/thebluefowl/skynet/db"
	"gorm.io/gorm"
)

type Metric struct {
	TimeToFirstByte        float64 `gorm:"column:time_to_first_byte"`
	FirstPaint             float64 `gorm:"column:first_paint" json:"first_paint"`
	FirstContentfulPaint   float64 `gorm:"column:first_contentful_paint" json:"first_contentful_paint"`
	FirstInputDelay        float64 `gorm:"column:first_input_delay" json:"first_input_delay"`
	LargestContentfulPaint float64 `gorm:"column:largest_contentful_paint" json:"largest_contentful_paint"`
	CumuLayoutShift        float64 `gorm:"column:cumu_layout_shift" json:"cumu_layout_shift"`
	TotalBlockingTime      float64 `gorm:"column:total_blocking_time" json:"total_blocking_time"`

	//Filters
	BrowserName         string `gorm:"type:varchar(50);column:browser_name" json:"browser_name"`
	BrowserVersion      string `gorm:"type:varchar(50);column:browser_version" json:"browser_version"`
	OperatingSystem     string `gorm:"type:varchar(50);column:operating_system" json:"operating_system"`
	NetworkInformation  string `gorm:"type:varchar(50);column:network_information" json:"network_information"`
	DeviceMemory        string `gorm:"type:varchar(50);column:device_memory" json:"device_memory"`
	HardwareConcurrency string `gorm:"type:varchar(50);column:hardware_concurrency" json:"hardware_concurrency"`
	ServiceWorkerStatus string `gorm:"type:varchar(50);column:service_worker_status" json:"service_worker_status"`
	IsLowEndDevice      bool   `gorm:"column:is_low_end_device" json:"is_low_end_device"`
	IsMobileDevice      bool   `gorm:"column:is_mobile_device" json:"is_mobile_device"`
	IsLowEndExperience  bool   `gorm:"column:is_low_end_experience" json:"is_low_end_experience"`

	// Timestamp
	Timestamp *time.Time `gorm:"column:timestamp;not null" json:"timestamp"`
}

type MetricStore struct {
	db   *sql.DB
	gorm *gorm.DB
}

func NewMetricStore() *MetricStore {
	_db, err := db.GetDB()
	if err != nil {
		log.Error("failed to get database connection, is it initialized?")
		panic(err)
	}
	conn, err := _db.GetSQLConn()
	if err != nil {
		log.Error("failed to get raw sql connection, is it initialized?")
		panic(err)
	}

	return &MetricStore{
		gorm: _db.DB,
		db:   conn,
	}
}

func (m *MetricStore) Create(metric *Metric) error {
	return m.gorm.Create(metric).Error
}

var MetricViewMap = map[string]string{
	"time_to_first_byte": "time_to_first_byte_hourly_percentiles",
	// "first_paint":              "first_paint_hourly_percentiles",
	// "first_contentful_paint":   "first_contentful_paint_hourly_percentiles",
	// "first_input_delay":        "first_input_delay_hourly_percentiles",
	// "largest_contentful_paint": "largest_contentful_paint_hourly_percentiles",
	// "cumu_layout_shift":        "cumu_layout_shift_hourly_percentiles",
	// "total_blocking_time":      "total_blocking_time_hourly_percentiles",
}

// Overview
const overviewBaseQuery = `SELECT
time_bucket_gapfill('1 day'::interval, bucket) as bucket,
locf(approx_percentile(0.75, rollup(pct_agg))) as p75
FROM %s
WHERE bucket >= now() - '30 day'::interval AND bucket < now()
GROUP BY 1
ORDER BY 1 DESC
LIMIT 2`

func generateOverviewQueryForView(view string) string {
	return fmt.Sprintf(overviewBaseQuery, view)
}

type MetricItem struct {
	Name       string  `json:"name"`
	Value      float64 `json:"value"`
	Change     float64 `json:"change"`
	IsPositive bool    `json:"is_positive"`
}

type OverviewResponse struct {
	Metrics []MetricItem `json:"metrics"`
}

func (m *MetricStore) GetOverview(field string) (*OverviewResponse, error) {
	//Get percentile value for time to first byte.
	overviewResponse := OverviewResponse{
		Metrics: []MetricItem{},
	}
	for metric_name, view := range MetricViewMap {
		query := generateOverviewQueryForView(view)
		rows, err := m.db.Query(query)
		if err != nil {
			log.Error("failed to query database")
			panic(err)
		}
		defer rows.Close()
		values := []float64{}
		i := 0
		for rows.Next() {
			if i > 1 {
				break
			}
			var bucket time.Time
			var p75 interface{}
			err := rows.Scan(&bucket, &p75)
			if err != nil {
				log.Error("failed to scan row")
				panic(err)
			}
			if p75 != nil {
				values = append(values, p75.(float64))
			} else {
				values = append(values, 0)
			}
			i++
		}
		if len(values) == 0 {
			values = []float64{0, 0}
		}
		if len(values) == 1 {
			values = []float64{values[0], 0}
		}
		responseItem := MetricItem{
			Name:       metric_name,
			Value:      values[0],
			Change:     values[0] - values[1],
			IsPositive: false,
		}
		overviewResponse.Metrics = append(overviewResponse.Metrics, responseItem)
		return &overviewResponse, nil
	}
	return nil, nil
}

// Trends
const trendBaseQuery = `SELECT
time_bucket_gapfill('1 day'::interval, bucket) as bucket,
locf(approx_percentile(0.75, rollup(pct_agg))) as p75
FROM %s
WHERE bucket >= now() - '30 day'::interval AND bucket < now()
GROUP BY 1
ORDER BY 1 DESC
LIMIT 30`

func generateTrendQueryForView(view string) string {
	return fmt.Sprintf(trendBaseQuery, view)
}

func (m *MetricStore) GetTrends(field string) (map[string][]float64, error) {
	//Get percentile value for time to first byte.
	trendResponse := map[string][]float64{}
	for metric_name, view := range MetricViewMap {
		query := generateTrendQueryForView(view)
		rows, err := m.db.Query(query)
		if err != nil {
			log.Error("failed to query database")
			panic(err)
		}
		defer rows.Close()
		values := []float64{}
		for rows.Next() {
			var bucket time.Time
			var p75 interface{}
			err := rows.Scan(&bucket, &p75)
			if err != nil {
				log.Error("failed to scan row")
				panic(err)
			}
			if p75 == nil {
				values = append(values, 0)
			} else {
				values = append(values, p75.(float64))
			}
		}
		trendResponse[metric_name] = values
	}
	return trendResponse, nil
}

// Web Stats

const webBrowserStatsQuery = `SELECT
time_bucket_gapfill('1 day'::interval, bucket) as bucket,
browser_name,
is_mobile_device,
locf(approx_percentile(0.75, rollup(pct_agg))) as p75
FROM %s
WHERE bucket >= now() - '30 day'::interval AND bucket < now() AND browser_name = '%s' AND is_mobile_device = false
GROUP BY 1, 2, 3
ORDER BY 1 DESC
LIMIT 2`

var Browsers = []string{"chrome"}

func (m *MetricStore) GetWebBrowserStats() (map[string][]MetricItem, error) {
	response := map[string][]MetricItem{}
	for _, browser := range Browsers {
		response[browser] = []MetricItem{}
		for metric_name, view := range MetricViewMap {
			query := fmt.Sprintf(webBrowserStatsQuery, view, browser)
			rows, err := m.db.Query(query)
			if err != nil {
				log.Error("failed to query database")
				panic(err)
			}
			defer rows.Close()

			values := []float64{}
			i := 0
			for rows.Next() {
				if i > 1 {
					break
				}
				i++
				var bucket time.Time
				var browser_name string
				var is_mobile_device bool
				var p75 interface{}
				err := rows.Scan(&bucket, &browser_name, &is_mobile_device, &p75)
				if err != nil {
					log.Error("failed to scan row")
					panic(err)
				}
				if p75 != nil {
					values = append(values, p75.(float64))
				} else {
					values = append(values, 0)
				}
			}
			if len(values) == 0 {
				values = []float64{0, 0}
			}
			if len(values) == 1 {
				values = []float64{values[0], 0}
			}

			responseItem := MetricItem{
				Name:       metric_name,
				Value:      values[0],
				Change:     values[0] - values[1],
				IsPositive: false,
			}
			response[browser] = append(response[browser], responseItem)
		}
	}
	return response, nil
}

// Mobile stats

const mobileStatsQuery = `SELECT
time_bucket_gapfill('1 day'::interval, bucket) as bucket,
browser_name,
is_mobile_device,
locf(approx_percentile(0.75, rollup(pct_agg))) as p75
FROM %s
WHERE bucket >= now() - '30 day'::interval AND bucket < now() AND browser_name = '%s' AND is_mobile_device = true
GROUP BY 1, 2, 3
ORDER BY 1 DESC
LIMIT 2`

var MobileBrowsers = []string{"chrome"}

func (m *MetricStore) GetMobileStats() (map[string][]MetricItem, error) {
	response := map[string][]MetricItem{}
	for _, browser := range MobileBrowsers {
		response[browser] = []MetricItem{}
		for metric_name, view := range MetricViewMap {
			query := fmt.Sprintf(mobileStatsQuery, view, browser)
			rows, err := m.db.Query(query)
			if err != nil {
				log.Error("failed to query database")
				panic(err)
			}
			defer rows.Close()

			values := []float64{}
			i := 0
			for rows.Next() {
				if i > 1 {
					break
				}
				var bucket time.Time
				var browser_name string
				var is_mobile_device bool
				var p75 interface{}
				err := rows.Scan(&bucket, &browser_name, &is_mobile_device, &p75)
				if err != nil {
					log.Error("failed to scan row")
					panic(err)
				}

				if p75 != nil {
					values = append(values, p75.(float64))
				} else {
					values = append(values, 0)
				}

			}
			if len(values) == 0 {
				values = []float64{0, 0}
			}
			if len(values) == 1 {
				values = []float64{values[0], 0}
			}

			responseItem := MetricItem{
				Name:       metric_name,
				Value:      values[0],
				Change:     values[0] - values[1],
				IsPositive: false,
			}
			response[browser] = append(response[browser], responseItem)
		}
	}
	return response, nil
}
