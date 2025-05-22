package usecase

import (
	"context"
	"errors"
	"fmt"
	"proyek1/internal/entity"
	"proyek1/internal/model"
	"proyek1/utils"
	"proyek1/utils/gmaps"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type RepositoryMapsInterface interface {
	InsertTempat(ctx context.Context, data *entity.Tempat) error
	GetTotalTempat(ctx context.Context, name string) (int, error)
	GetTempatPagination(ctx context.Context, name string, limit, offset int) ([]entity.Tempat, error)
	GetDetailTempat(ctx context.Context, id string) (entity.GetDetailTempat, error)
}

type UsecaseMaps struct {
	repo RepositoryMapsInterface
	gm   gmaps.GmapsInterface
	log  *logrus.Logger
}

func NewMapsUsercase(repo RepositoryMapsInterface, log *logrus.Logger, gm gmaps.GmapsInterface) *UsecaseMaps {
	return &UsecaseMaps{
		repo: repo,
		log:  log,
		gm:   gm,
	}
}

func (s *UsecaseMaps) InsertTempat(ctx context.Context, placeId string) error {
	if placeId == "" {
		return errors.New("Id tidak ditemukan atau kosong")
	}
	dataGmaps, err := s.gm.GmapsSearchByPlaceID(placeId)
	if err != nil {
		return err
	}

	err = s.repo.InsertTempat(ctx, ConverMapsToModelPlace(dataGmaps))
	if err != nil {
		return err
	}

	return nil
}
func (s *UsecaseMaps) GetTempatPagination(ctx context.Context, name string, limit, page int) ([]model.GetAllTempat, int, error) {
	total, err := s.repo.GetTotalTempat(ctx, name)
	if err != nil {
		return []model.GetAllTempat{}, 0, err
	}
	totalPage := utils.TotalPageForPagination(total, limit)
	offset := (page - 1) * limit
	dataTempat, err := s.repo.GetTempatPagination(ctx, name, limit, offset)
	if err != nil {
		return []model.GetAllTempat{}, 0, err
	}

	var res []model.GetAllTempat
	for _, v := range dataTempat {
		tempat := model.GetAllTempat{
			ID:      v.ID,
			PlaceId: v.PlaceId,
			Name:    v.Name,
			Address: v.Address,
		}

		var hours []model.HourTempatGetAll
		for _, h := range v.OpeningHours {
			hours = append(hours, model.HourTempatGetAll{
				Day:       h.Day,
				OpenTime:  h.OpenTime,
				CloseTime: h.CloseTime,
			})
		}

		var foto []model.FotoTempatGetAll
		for _, f := range v.Photos {
			// _, err := s.gm.PhotoReference(f.PhotoRefrences)
			// if err != nil {
			// 	s.log.Warn("Photo reference error: ", err)
			// 	continue
			// }

			proxyURL := f.PhotoRefrences
			foto = append(foto, model.FotoTempatGetAll{
				WidthPx:        f.WidthPx,
				HeightPx:       f.HeightPx,
				PhotoRefrences: proxyURL,
			})
		}

		tempat.OpeningHours = hours
		tempat.Photos = foto

		res = append(res, tempat)

	}

	return res, totalPage, nil
}

func (s *UsecaseMaps) GetDetailTempat(ctx context.Context, id string) (model.GetDetailTempat, error) {
	if id == "" {
		return model.GetDetailTempat{}, errors.New("Id tidak ditemukan atau kosong")
	}

	resData, err := s.repo.GetDetailTempat(ctx, id)
	if err != nil {
		return model.GetDetailTempat{}, err
	}
	var totalRatingSum float64 = 0.0
	var parsedReviews []model.Review
	for _, r := range resData.Reviews {
		totalRatingSum += float64(r.Rating)
		parsedReviews = append(parsedReviews, model.Review{
			AuthorName:                     r.Author,
			Text:                           r.Text,
			Rating:                         float64(r.Rating),
			RelativePublishTimeDescription: r.ReviewCreated,
		})
	}

	var averageRating float64
	if len(resData.Reviews) > 0 {
		averageRating = totalRatingSum / float64(len(resData.Reviews))
	} else {
		averageRating = 0.0
	}
	var periods []model.HourTempatGetAll
	for _, p := range resData.OpeningHours {
		periods = append(periods, model.HourTempatGetAll{
			Day:       p.Day,
			OpenTime:  p.OpenTime,
			CloseTime: p.CloseTime,
		})
	}
	var photos []model.Photo
	for _, s := range resData.Photos {
		photos = append(photos, model.Photo{
			WidthPx:        s.WidthPx,
			HeightPx:       s.HeightPx,
			PhotoReference: s.PhotoRefrences,
		})
	}
	var types []model.Type
	for _, s := range resData.Types {
		types = append(types, model.Type{
			CategoryCode: s.CategoryCode,
			PlaceID:      s.PlaceID,
		})
	}
	results := model.GetDetailTempat{
		PlaceID:          resData.PlaceID,
		Name:             resData.Name,
		FormattedAddress: resData.FormattedAddress,
		Lat:              fmt.Sprintf("%f", resData.Lat),
		Lang:             fmt.Sprintf("%f", resData.Lng),
		Icon:             resData.Icon,
		Rating:           averageRating,
		Reviews:          parsedReviews,
		RegularOpeningHours: model.DetailHour{
			Periods: periods,
		},
		NavigasiURL: fmt.Sprintf("https://www.google.com/maps/search/?api=1&query=%f,%f&query_place_id=%s",
			resData.Lat,
			resData.Lng,
			id),
		Photos:         photos,
		Types:          types,
		BusinessStatus: resData.BusinessStatus,
	}

	return results, nil
}

func (s *UsecaseMaps) RouteDestination(ctx context.Context, req model.RequestRouteMaps, placeID string) (*model.ResponseRouteMaps, error) {
	searchData, err := s.gm.GmapsSearchByPlaceID(placeID)
	if err != nil {
		return nil, err
	}

	floatLat, _ := strconv.ParseFloat(searchData.Geometry.Lat, 64)
	floatLng, _ := strconv.ParseFloat(searchData.Geometry.Lng, 64)
	reqData := model.RequestRouteMaps{
		Origin: model.Waypoint{
			Location: model.LocationReq{
				LatLng: model.LatLng{
					Latitude:  req.Origin.Location.LatLng.Latitude,
					Longitude: req.Origin.Location.LatLng.Longitude,
				},
			},
		}, Destination: model.Waypoint{
			Location: model.LocationReq{
				LatLng: model.LatLng{
					Latitude:  floatLat,
					Longitude: floatLng,
				},
			},
		},
		TravelMode: req.TravelMode,
	}
	fmt.Println("Hasil pencarian placeID:", searchData.Geometry.Lat, searchData.Geometry.Lng)

	return s.gm.RouteToDestination(reqData)
}

// Convert
func ConverMapsToModelPlace(req model.MapsGetByPlaceId) *entity.Tempat {
	lat, _ := strconv.ParseFloat(req.Geometry.Lat, 64)
	lng, _ := strconv.ParseFloat(req.Geometry.Lng, 64)
	conv := &entity.Tempat{
		ID:             uuid.New().String(),
		PlaceId:        req.PlaceID,
		Name:           req.Name,
		Latitude:       lat,
		Longtitude:     lng,
		Address:        req.FormattedAddress,
		Icon:           req.Icon,
		BusinessStatus: req.BusinessStatus,
	}

	var rev []entity.Review
	for _, v := range req.Reviews {
		rev = append(rev, entity.Review{
			ID:            uuid.New().String(),
			PlaceId:       req.PlaceID,
			UserId:        nil,
			Author:        v.AuthorName,
			ReviewCreated: v.RelativePublishTimeDescription,
			Text:          v.Text,
			Rating:        int(v.Rating),
			IsFromGoogle:  true,
		})
	}

	var photos []entity.Photo
	for _, p := range req.Photos {
		photos = append(photos, entity.Photo{
			ID:             uuid.New().String(),
			PlaceId:        req.PlaceID,
			UserId:         nil,
			ReviewID:       nil,
			PhotoRefrences: p.PhotoReference,
			WidthPx:        p.WidthPx,
			HeightPx:       p.HeightPx,
			IsFromGoogle:   true,
		})
	}

	var hours []entity.Hour
	for _, j := range req.RegularOpeningHours.Periods {
		hours = append(hours, entity.Hour{
			ID:        uuid.New().String(),
			PlaceId:   req.PlaceID,
			Day:       strconv.Itoa(j.Open.Day),
			OpenTime:  j.Open.Time,
			CloseTime: j.Close.Time,
		})
	}

	var types []entity.Type
	for _, t := range req.Types {
		types = append(types, entity.Type{
			CategoryCode: t,
			PlaceID:      req.PlaceID,
		})
	}

	conv.Types = types
	conv.Reviews = rev
	conv.Photos = photos
	conv.OpeningHours = hours

	return conv
}
