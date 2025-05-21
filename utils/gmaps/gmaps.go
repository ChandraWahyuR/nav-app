package gmaps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"proyek1/config"
	"proyek1/constant"
	"proyek1/internal/model"
	"proyek1/utils"
)

type GmapsInterface interface {
	GmapsSearchObject(inputTempat string) (model.Maps, error)
	GmapsSearchList(inputTempat string) ([]model.Maps, error)
	GmapsSearchByPlaceID(placeID string) (model.MapsGetByPlaceId, error)
	PhotoReference(photoURl string) (string, error)
	RouteToDestination(req model.RequestRouteMaps) (*model.ResponseRouteMaps, error)
}

type gmapsStruct struct {
	c config.GMAPS
}

func NewMail(c config.GMAPS) gmapsStruct {
	return gmapsStruct{
		c: c,
	}
}

func (c *gmapsStruct) GmapsSearchObject(inputTempat string) (model.Maps, error) {
	encodedInput := url.QueryEscape(inputTempat)
	requestURL := fmt.Sprintf("%s%s&inputtype=textquery&key=%s", constant.Gmaps, encodedInput, c.c.GMAPS_API_KEY)

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error melakukan permintaan:", err)
		return model.Maps{}, err
	}
	defer resp.Body.Close()

	// Baca response dari api
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error membaca body respons:", err)
		return model.Maps{}, err
	}

	var searchResponse model.GmapsAPIGetObject
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return model.Maps{}, fmt.Errorf("error unmarshal: %w", err)
	}

	var results model.Maps
	for _, v := range searchResponse.Candidates {
		results = model.Maps{
			PlaceID: v.PlaceID,
			Name:    v.Name,
			Geometry: model.LocationResp{
				Lat: fmt.Sprintf(`%f`, v.Geometry.Location.Lat),
				Lng: fmt.Sprintf(`%f`, v.Geometry.Location.Lng),
			},
		}
	}

	return results, nil
}

func (c *gmapsStruct) GmapsSearchList(inputTempat string) ([]model.Maps, error) {
	encodedInput := url.QueryEscape(inputTempat)
	requestURL := fmt.Sprintf("%s?query=%s&key=%s", constant.GmapsSearchText, encodedInput, c.c.GMAPS_API_KEY)

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error saat request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error membaca response body:", err)
		return nil, err
	}

	var searchResponse model.GmapsAPIGetTextSearch
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return nil, fmt.Errorf("error unmarshal response: %w", err)
	}

	var results []model.Maps
	for _, v := range searchResponse.Place {
		results = append(results, model.Maps{
			PlaceID: v.PlaceID,
			Name:    v.Name,
			Geometry: model.LocationResp{
				Lat: fmt.Sprintf(`%f`, v.Geometry.Location.Lat),
				Lng: fmt.Sprintf(`%f`, v.Geometry.Location.Lng),
			},
		})
	}

	return results, nil
}

func (c *gmapsStruct) GmapsSearchByPlaceID(placeID string) (model.MapsGetByPlaceId, error) {
	encodedInput := url.QueryEscape(placeID)
	requestURL := fmt.Sprintf("%s=%s&language=id&key=%s", constant.GmapsGetByPlaceID, encodedInput, c.c.GMAPS_API_KEY)

	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Println("Error saat request:", err)
		return model.MapsGetByPlaceId{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error membaca response body:", err)
		return model.MapsGetByPlaceId{}, err
	}

	var searchResponse model.GmapsAPIGetPlaceDetails
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		return model.MapsGetByPlaceId{}, fmt.Errorf("error unmarshal response: %w", err)
	}

	var parsedReviews []model.Review
	for _, r := range searchResponse.Place.Reviews {
		parsedReviews = append(parsedReviews, model.Review{
			AuthorName:                     r.AuthorName,
			Text:                           r.Text,
			Rating:                         r.Rating,
			RelativePublishTimeDescription: r.RelativePublishTimeDescription,
		})
	}

	var periods []model.Period
	for _, p := range searchResponse.Place.RegularOpeningHours.Periods {
		periods = append(periods, model.Period{
			Open: model.DayTime{
				Day:  p.Open.Day,
				Time: utils.FormatJam(p.Open.Time),
			},
			Close: model.DayTime{
				Day:  p.Close.Day,
				Time: utils.FormatJam(p.Close.Time),
			},
		})
	}
	var photos []model.Photo
	for _, s := range searchResponse.Place.Photos {
		photos = append(photos, model.Photo{
			WidthPx:        s.WidthPx,
			HeightPx:       s.HeightPx,
			PhotoReference: s.PhotoReference,
		})
	}
	results := model.MapsGetByPlaceId{
		PlaceID:          searchResponse.Place.PlaceID,
		Name:             searchResponse.Place.Name,
		FormattedAddress: searchResponse.Place.FormattedAddress,
		Geometry: model.LocationResp{
			Lat: fmt.Sprintf("%f", searchResponse.Place.Geometry.Location.Lat),
			Lng: fmt.Sprintf("%f", searchResponse.Place.Geometry.Location.Lng),
		},
		Icon:    searchResponse.Place.Icon,
		Rating:  searchResponse.Place.Rating,
		Reviews: parsedReviews,
		RegularOpeningHours: model.OpeningHour{
			OpenNow: searchResponse.Place.RegularOpeningHours.OpenNow,
			Periods: periods,
		},
		NavigasiURL: fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f&query_place_id=%s",
			searchResponse.Place.Geometry.Location.Lat,
			searchResponse.Place.Geometry.Location.Lng,
			placeID),
		Photos:         photos,
		BusinessStatus: searchResponse.Place.BusinessStatus,
		Types:          searchResponse.Place.Types,
	}

	return results, nil
}

func (c *gmapsStruct) PhotoReference(photoURl string) (string, error) {
	if photoURl == "" {
		return "", fmt.Errorf("empty photo reference")
	}
	photoURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/photo?maxwidth=400&photo_reference=%s&key=%s",
		photoURl,
		c.c.GMAPS_API_KEY,
	)

	return photoURL, nil
}

func (c *gmapsStruct) RouteToDestination(req model.RequestRouteMaps) (*model.ResponseRouteMaps, error) {
	var client = &http.Client{}

	requestURL := fmt.Sprintf("%s?key=%s", constant.GmapsGetRouteByPlaceID, c.c.GMAPS_API_KEY)
	fmt.Println(requestURL)
	jsonData, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", requestURL, bytes.NewReader(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// fieldMask := "routes.duration,routes.distanceMeters,routes.polyline.encodedPolyline"
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Goog-FieldMask", "*")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	fmt.Println("ini response.body:", response.Body)
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)
	fmt.Println("Full Response Body:\n", bodyString) // Cetak seluruh body

	res := &model.ResponseRouteMaps{}
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		return nil, err
	}
	fmt.Println("ini res:", res)

	return res, nil
}
