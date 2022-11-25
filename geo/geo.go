package geo

import (
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
