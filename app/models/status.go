package models

type Status struct {
	IP       string `json:"ip"`
	State    string `json:"state"`
	SubState string `json:"subState"`
}
