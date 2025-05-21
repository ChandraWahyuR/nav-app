package constant

const (
	Gmaps                  = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json?fields=place_id,name,geometry&input=" // findplcaefromtext = return object, jadi pakai textsearch  kalau array
	GmapsSearchText        = "https://maps.googleapis.com/maps/api/place/textsearch/json"                                             // findplcaefromtext = return object, jadi pakai textsearch  kalau array
	GmapsGetByPlaceID      = "https://maps.googleapis.com/maps/api/place/details/json?place_id"
	GmapsGetRouteByPlaceID = "https://routes.googleapis.com/directions/v2:computeRoutes"
	VercelRoute            = "https://html-411k7ckwk-chands-projects-5f68fc9c.vercel.app/static/index.html"

	// Message Response
	StatusSuccess = "Success"
	StatusFail    = "Fail"
)
