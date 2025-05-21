package model

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

// Get Detailed
type GetDetailTempat struct {
	PlaceID             string `json:"place_id"`
	Name                string `json:"name"`
	Lat                 string
	Lang                string
	FormattedAddress    string     `json:"formatted_address"`
	Icon                string     `json:"icon"`
	NavigasiURL         string     `json:"navigasi_url"`
	Rating              float64    `json:"rating"`
	Reviews             []Review   `json:"reviews"`
	RegularOpeningHours DetailHour `json:"current_opening_hours"`
	Photos              []Photo    `json:"photos"`
	BusinessStatus      string     `json:"business_status"`
	Types               []string   `json:"types"`
}

type DetailHour struct {
	Periods []HourTempatGetAll `json:"periods"`
}
