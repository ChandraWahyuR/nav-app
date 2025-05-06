package model

// Object
type GmapsAPIGetObject struct {
	Candidates []struct {
		PlaceID  string   `json:"place_id"`
		Name     string   `json:"name"`
		Geometry Geometry `json:"geometry"`
	} `json:"candidates"`
	Status string `json:"status"`
}

// List
type GmapsAPIGetTextSearch struct {
	HTMLAttributions []interface{} `json:"html_attributions"`
	Place            []PlaceResult `json:"results"`
	NextPageToken    string        `json:"next_page_token"`
	Status           string        `json:"status"`
}

type GmapsAPIGetPlaceDetails struct {
	Place  PlaceResult `json:"result"`
	Status string      `json:"status"`
}

//
type PlaceResult struct {
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

type Photo struct {
	WidthPx        int    `json:"width"`
	HeightPx       int    `json:"height"`
	PhotoReference string `json:"photo_reference"`
}

type AuthorAttribution struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type Review struct {
	AuthorName                     string  `json:"author_name"`
	RelativePublishTimeDescription string  `json:"relative_time_description"`
	Text                           string  `json:"text"`
	Rating                         float64 `json:"rating"`
}
