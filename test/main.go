package main

import (
	"fmt"

	"github.com/thebluefowl/skynet/geo"
)

func main() {
	geo.LoadCountryToLatLOng()
	country, _ := geo.GetCountry("203.192.247.242")
	fmt.Println(geo.CountryLatLong[country])
}
