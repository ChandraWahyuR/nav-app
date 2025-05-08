package entity

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
