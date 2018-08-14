package cmd

import (
	"testing"

	"github.com/breiting/rex"
)

func TestGetGeoLocation(t *testing.T) {
	var address rex.ProjectAddress

	address.AddressLine1 = "Stremayrgasse 16"

	lat, lon := GetGeoLocation(&address)

	if lat != 47.0650332 {
		t.Errorf("Latitude is wrong")
	}
	if lon != 15.4522326353324 {
		t.Errorf("Longitude is wrong")
	}
}
