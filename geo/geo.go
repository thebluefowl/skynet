package geo

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/ip2location/ip2location-go/v9"
)

type LatLong struct {
	Lat  float64
	Long float64
}

var CountryLatLong = map[string]LatLong{}

func GetCountry(ip string) (string, float64, float64, error) {
	db, err := ip2location.OpenDB("IP2LOCATION-LITE-DB5.IPV6.BIN")
	if err != nil {
		return "", 0, 0, err
	}

	results, err := db.Get_all(ip)
	if err != nil {
		return "", 0, 0, err
	}
	return results.Country_long, float64(results.Latitude), float64(results.Longitude), nil
}

func LoadCountryToLatLOng() error {
	f, err := os.Open("country-lat-long.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = ','

	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		lat, _ := strconv.ParseFloat(record[1], 64)
		long, _ := strconv.ParseFloat(record[2], 64)
		CountryLatLong[record[0]] = LatLong{
			Lat:  lat,
			Long: long,
		}
	}
	return nil
}
