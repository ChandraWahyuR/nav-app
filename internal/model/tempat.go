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

// Get
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
