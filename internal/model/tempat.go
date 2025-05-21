package model

type Tempat struct {
	ID             string
	PlaceId        string
	Name           string
	Latitude       float64
	Longtitude     float64
	Address        string
	Icon           string
	BusinessStatus string

	Reviews      []ReviewTempat
	Photos       []PhotoTempat
	OpeningHours []Hour
}

type ReviewTempat struct {
	ID            string
	PlaceId       string
	UserId        string
	Author        string
	ReviewCreated string
	Text          string
	Address       string
	Rating        int
	IsFromGoogle  bool
	Photos        []Photo
}

type PhotoTempat struct {
	ID             string
	PlaceId        string
	UserId         string
	ReviewID       string
	PhotoRefrences string `json:"photo_reference"`
	WidthPx        int    `json:"width_px"`
	HeightPx       int    `json:"height_px"`
	IsFromGoogle   bool
}

type Hour struct {
	ID        string
	PlaceId   string
	Day       string
	OpenTime  string
	CloseTime string
}

// ==========================================================================================
// Get All
type GetAllTempat struct {
	ID           string             `json:"id"`
	PlaceId      string             `json:"place_id"`
	Name         string             `json:"name"`
	Address      string             `json:"address"`
	Photos       []FotoTempatGetAll `json:"photos"`
	OpeningHours []HourTempatGetAll `json:"opening_hours"`
}

type FotoTempatGetAll struct {
	PhotoRefrences string `json:"photo_reference"`
	WidthPx        int    `json:"width_px"`
	HeightPx       int    `json:"height_px"`
}

type HourTempatGetAll struct {
	Day       string `json:"day"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

type GetDetailTempat struct {
	PlaceID             string `json:"place_id"`
	Name                string `json:"name"`
	Lat                 string
	Lang                string
	FormattedAddress    string      `json:"formatted_address"`
	Icon                string      `json:"icon"`
	NavigasiURL         string      `json:"navigasi_url"`
	Rating              float64     `json:"rating"`
	Reviews             []Review    `json:"reviews"`
	RegularOpeningHours OpeningHour `json:"current_opening_hours"`
	Photos              []Photo     `json:"photos"`
	BusinessStatus      string      `json:"business_status"`
	Types               []string    `json:"types"`
}
