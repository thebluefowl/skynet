package main

import (
	"fmt"
	"log"

	"github.com/thebluefowl/skynet/db"
	"github.com/thebluefowl/skynet/model"
)

const QueryCreateHypertable = `SELECT create_hypertable('metrics', 'timestamp');`
const QueryCreateMaterializedViewBase = `CREATE MATERIALIZED VIEW %s
WITH (timescaledb.continuous)
AS SELECT
    time_bucket('1 day'::interval, timestamp) as bucket,
    percentile_agg(%s) as pct_agg,
	browser_name,
	is_mobile_device
FROM metrics
GROUP BY 1, browser_name, is_mobile_device;
`

func CreateTables() {
	db, err := db.GetDB()
	if err != nil {
		panic(err)
	}

	db.DB.AutoMigrate(&model.Metric{})
	conn, err := db.GetSQLConn()
	if err != nil {
		panic(err)
	}
	log.Println("Creating hypertable")
	_, err = conn.Exec(QueryCreateHypertable)
	if err != nil {
		panic(err)
	}

	for metric_name, view := range model.MetricViewMap {
		log.Printf("creating materialized view")
		_, err = conn.Exec(fmt.Sprintf(QueryCreateMaterializedViewBase, view, metric_name))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	CreateTables()
}

/*
CREATE MATERIALIZED VIEW foo_hourly
WITH (timescaledb.continuous)
AS SELECT
    time_bucket('1 h'::interval, time) as bucket,
    percentile_agg(first_paint) as pct_agg
FROM metrics
GROUP BY 1;

SELECT
    time_bucket('1 h'::interval, bucket) as bucket,
    approx_percentile(0.75, rollup(pct_agg)) as p75,
    approx_percentile(0.95, rollup(pct_agg)) as p95
FROM metrics_hourly_percentiles
GROUP BY 1;

INSERT INTO metrics VALUES (1,1,2,36,2,3,4, 'chrome', '3.2.3', 'linux', '4g', '8g', '4', 'active', true, false, NOW()-'1 day'::interval);
INSERT INTO metrics VALUES (2,3,1,2,36,2, 'chrome', '2.2.3', 'linux', '4g', '8g', '4', 'active', true, false, DATE_SUB(NOW(), INTERVAL '1' HOUR));


SELECT browser_name, histogram(first_paint, 1, 2, 5)
FROM metrics
GROUP BY browser_name
LIMIT 10;

SELECT
    time_bucket('1 s'::interval, time),
    average(stats_agg(first_paint))
FROM metrics
GROUP BY 1;


CREATE MATERIALIZED VIEW foo_hourly
WITH (timescaledb.continuous)
AS SELECT
    time_bucket('1 min'::interval, time) as bucket,
    percentile_agg(first_paint) as pct_agg,
FROM metrics
GROUP BY 1, browser_name;


SELECT
time_bucket_gapfill('1 h'::interval, bucket) as bucket,
locf(approx_percentile(0.75, rollup(pct_agg))) as p75,
locf(approx_percentile(0.95, rollup(pct_agg))) as p95
FROM time_to_first_byte_hourly_percentiles
WHERE bucket >= now() - '1 day'::interval AND bucket < now()
GROUP BY 1
ORDER BY 1 DESC
LIMIT 2


SELECT
time_bucket_gapfill('1 day'::interval, bucket) as bucket,
locf(approx_percentile(0.75, rollup(pct_agg))) as p75
FROM
WHERE bucket >= now() - '30 day'::interval AND bucket < now()
GROUP BY 1
ORDER BY 1 DESC
LIMIT 2
*/
