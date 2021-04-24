package geoip

type GeoInfo struct {
	Coordinates `json:"coordinates"`
	Country     `json:"country"`
	Region      `json:"region"`
	City        `json:"city"`
}

type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Region struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type City struct {
	Name string `json:"name"`
}
