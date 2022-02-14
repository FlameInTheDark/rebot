package geonames

import "strconv"

// LocationResult contains found locations
type LocationResult struct {
	ResultCount int       `json:"totalResultsCount"`
	Geonames    []Geoname `json:"geonames"`
}

// Geoname contains location data
type Geoname struct {
	GeonameID   int    `json:"geonameId"`
	CountryID   string `json:"countryId"`
	ToponymName string `json:"toponymName"`
	Population  int    `json:"population"`
	CountryCode string `json:"countryCode"`
	Name        string `json:"name"`
	CountryName string `json:"countryName"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
}

// Coordinates returns latitude and longitude
func (g *Geoname) Coordinates() (string, string) {
	return g.Lat, g.Lng
}

// CoordinatesFloat64 converts coordinates to float64
func (g *Geoname) CoordinatesFloat64() (float64, float64, error) {
	lat, err := strconv.ParseFloat(g.Lat, 64)
	if err != nil {
		return 0, 0, err
	}
	lng, err := strconv.ParseFloat(g.Lng, 64)
	if err != nil {
		return 0, 0, err
	}
	return lat, lng, nil
}
