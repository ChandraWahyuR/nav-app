package constant

const (
	Gmaps             = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json?fields=place_id,name,geometry&input=" // findplcaefromtext = return object, jadi pakai textsearch  kalau array
	GmapsSearchText   = "https://maps.googleapis.com/maps/api/place/textsearch/json"                                             // findplcaefromtext = return object, jadi pakai textsearch  kalau array
	GmapsGetByPlaceID = "https://maps.googleapis.com/maps/api/place/details/json?place_id"
)
