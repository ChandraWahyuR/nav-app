package entity

// Get All
type Tempat struct {
	ID             string
	PlaceId        string
	Name           string
	Latitude       float64
	Longtitude     float64
	Address        string
	Icon           string
	BusinessStatus string

	Reviews      []Review
	Photos       []Photo
	OpeningHours []Hour
}

type Review struct {
	ID            string
	PlaceId       string
	UserId        *string
	Author        string
	ReviewCreated string
	Text          string
	Rating        int
	IsFromGoogle  bool
	Photos        []Photo
}

type Photo struct {
	ID             string
	PlaceId        string
	UserId         *string
	ReviewID       *string
	PhotoRefrences string `json:"photo_reference"`
	WidthPx        int    `json:"width_px"`
	HeightPx       int    `json:"height_px"` // buat unmarshal
	IsFromGoogle   bool
}

type Hour struct {
	ID        string `json:"id"`
	PlaceId   string `json:"place_id"`
	Day       string `json:"day"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
}

// ==========================================================================================================================
// Get Detail
type GetDetailTempat struct {
	PlaceID             string      `json:"place_id"`
	Name                string      `json:"name"`
	FormattedAddress    string      `json:"formatted_address"`
	NavigasiURL         string      `json:"navigasi_url"`
	Geometry            Geometry    `json:"geometry"`
	Icon                string      `json:"icon"`
	Rating              float64     `json:"rating"`
	Reviews             []Review    `json:"reviews"`
	RegularOpeningHours OpeningHour `json:"current_opening_hours"`
	Photos              []Photo     `json:"photos"`
	BusinessStatus      string      `json:"business_status"`
	Types               []string    `json:"types"`
}

type Geometry struct {
	Location Location
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type OpeningHour struct {
	OpenNow bool     `json:"open_now"`
	Periods []Period `json:"periods"`
}

type Period struct {
	Open  DayTime `json:"open"`
	Close DayTime `json:"close"`
}

type DayTime struct {
	Day  int    `json:"day"`
	Time string `json:"time"`
}

//
