package model

// Response Object
type GmapsResponse struct {
	Candidates []Maps `json:"candidates"`
	Status     string `json:"status"`
}

type Maps struct {
	PlaceID  string       `json:"place_id"`
	Name     string       `json:"name"`
	Geometry LocationResp `json:"geometry"`
}

type LocationResp struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type MapsGetByPlaceId struct {
	PlaceID             string       `json:"place_id"`
	Name                string       `json:"name"`
	Geometry            LocationResp `json:"geometry"`
	FormattedAddress    string       `json:"formatted_address"`
	Icon                string       `json:"icon"`
	NavigasiURL         string       `json:"navigasi_url"`
	Rating              float64      `json:"rating"`
	Reviews             []Review     `json:"reviews"`
	RegularOpeningHours OpeningHour  `json:"current_opening_hours"`
	Photos              []Photo      `json:"photos"`
	BusinessStatus      string       `json:"business_status"`
	Types               []string     `json:"types"`
}
