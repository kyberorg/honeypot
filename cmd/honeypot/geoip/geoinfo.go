package geoip

type GeoInfo struct {
	Location `json:"location,omitempty"`
	Country  `json:"country,omitempty"`
	Region   `json:"region,omitempty"`
	City     `json:"city,omitempty"`
}

type Location struct {
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
