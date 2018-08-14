/*
 * Author: Bernhard Reitinger
 * Date  : 2018
 */

package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/breiting/rex"
	"github.com/tidwall/gjson"
)

func console(err error, value interface{}) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", value)
	}
}

// GetGeoLocation gets the best guess geo-location from OpenStreetMap
func GetGeoLocation(address *rex.ProjectAddress) (lat, lon float64) {

	// Spaces are not allowed
	street := strings.Replace(address.AddressLine1, " ", "+", -1)
	postcode := strings.Replace(address.PostCode, " ", "+", -1)
	city := strings.Replace(address.City, " ", "+", -1)

	req := fmt.Sprintf("%s?format=json", openStreetMapSearch)
	req += fmt.Sprintf("&street=%s", street)
	req += fmt.Sprintf("&postalcode=%s", postcode)
	req += fmt.Sprintf("&city=%s", city)
	req += fmt.Sprintf("&limit=1")

	// fmt.Println("Getting OpenStreetMap data: ", req)
	r, err := http.Get(req)
	if err != nil {
		return 0.0, 0.0
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	b := string(body)
	// fmt.Println("Response: ", b)

	lat = gjson.Get(string(body), "0.lat").Float()
	lon = gjson.Get(string(body), "0.lon").Float()

	fmt.Printf("OpenStreetMap data\n\n")
	fmt.Println("LICENSE notice:   ", gjson.Get(b, "0.licence").String())
	fmt.Println("Geolocation name: ", gjson.Get(b, "0.display_name").String())
	fmt.Println("Geolocation lat:  ", lat)
	fmt.Println("Geolocation lon:  ", lon)
	fmt.Printf("\n")

	return lat, lon
}
