package dto

type CollectedData struct {
	Time     string `json:"time"`
	User     string `json:"user"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}
