package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"proyek1/internal/entity"
	"proyek1/utils"

	"github.com/sirupsen/logrus"
)

type MapsRepo struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewMapsRepository(db *sql.DB, log *logrus.Logger) *MapsRepo {
	return &MapsRepo{
		db:  db,
		log: log,
	}
}

func (r *MapsRepo) GetTotalTempat(ctx context.Context) (int, error) {
	var total int
	query := `
		SELECT COUNT(DISTINCT tempat_pariwisata.place_id)
		FROM tempat_pariwisata
		INNER JOIN foto_tempat ON foto_tempat.place_id = tempat_pariwisata.place_id
		WHERE tempat_pariwisata.deleted_at IS NULL;
	`
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *MapsRepo) GetTempatPagination(ctx context.Context, limit, offset int) ([]entity.Tempat, error) {
	var res []entity.Tempat

	query := `SELECT tempat_pariwisata.id, tempat_pariwisata.place_id, tempat_pariwisata.name, tempat_pariwisata.address, tempat_pariwisata.icon, 
				json_agg(
					json_build_object(
						'photo_reference', foto_tempat.photo_reference,
						'width_px', CAST(foto_tempat.width_px AS INTEGER),
						'height_px', CAST(foto_tempat.height_px AS INTEGER)
					)
				) AS photos,
				json_agg(
					json_build_object(
						'day', opening_hours.day,
						'open_time', opening_hours.open_time,
						'close_time', opening_hours.close_time
					)
				) AS time
				FROM "tempat_pariwisata" 
				INNER JOIN foto_tempat 
					ON foto_tempat.place_id = tempat_pariwisata.place_id 
				LEFT JOIN opening_hours 
    				ON opening_hours.place_id = tempat_pariwisata.place_id
				WHERE tempat_pariwisata.deleted_at IS NULL 
					GROUP BY 
					tempat_pariwisata.id, tempat_pariwisata.place_id, tempat_pariwisata.name, tempat_pariwisata.address, tempat_pariwisata.icon
					LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return []entity.Tempat{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var photoJson, timeJson []byte

		var tempat entity.Tempat
		err := rows.Scan(&tempat.ID, &tempat.PlaceId, &tempat.Name, &tempat.Address, &tempat.Icon, &photoJson, &timeJson)
		if err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}

		//
		if err := json.Unmarshal(photoJson, &tempat.Photos); err != nil {
			return nil, fmt.Errorf("error unmarshalling photo_json: %w", err)
		}

		if err := json.Unmarshal(timeJson, &tempat.OpeningHours); err != nil {
			return nil, fmt.Errorf("error unmarshalling time_json: %w", err)
		}
		res = append(res, tempat)
	}

	return res, nil
}

func (r *MapsRepo) InsertTempat(ctx context.Context, data *entity.Tempat) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert tempat
	query := `INSERT INTO tempat_pariwisata (id, place_id, name, latitude, longtitude, address, icon, business_status)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.ExecContext(ctx, query, data.ID, data.PlaceId, data.Name, data.Latitude, data.Longtitude, data.Address, data.Icon, data.BusinessStatus)
	if err != nil {
		return utils.ParsePQError(err)
	}

	//
	if len(data.Reviews) > 0 {
		if err := r.InsertReview(ctx, tx, data.Reviews); err != nil {
			return utils.ParsePQError(err)
		}
	}
	if len(data.Photos) > 0 {
		if err := r.InsertPhotos(ctx, tx, data.Photos); err != nil {
			return utils.ParsePQError(err)
		}
	}
	if len(data.OpeningHours) > 0 {
		if err := r.InsertHours(ctx, tx, data.OpeningHours); err != nil {
			return utils.ParsePQError(err)
		}
	}

	return tx.Commit()
}

func (r *MapsRepo) InsertReview(ctx context.Context, tx *sql.Tx, data []entity.Review) error {
	for _, review := range data {
		q := `INSERT INTO review_tempat (id, place_id, users_id, author, review_created, text, rating, isfrom_google)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := tx.ExecContext(ctx, q, review.ID, review.PlaceId, review.UserId, review.Author, review.ReviewCreated, review.Text, review.Rating, review.IsFromGoogle)
		if err != nil {
			return utils.ParsePQError(err)
		}
	}

	return nil
}

func (r *MapsRepo) InsertPhotos(ctx context.Context, tx *sql.Tx, data []entity.Photo) error {
	for _, photo := range data {
		q := `INSERT INTO foto_tempat (id, place_id, review_id, users_id, photo_reference, width_px, height_px, isfrom_google)
				  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err := tx.ExecContext(ctx, q, photo.ID, photo.PlaceId, photo.ReviewID, photo.UserId, photo.PhotoRefrences, photo.WidthPx, photo.HeightPx, photo.IsFromGoogle)
		if err != nil {
			return utils.ParsePQError(err)
		}
	}
	return nil
}

func (r *MapsRepo) InsertHours(ctx context.Context, tx *sql.Tx, data []entity.Hour) error {
	for _, hour := range data {
		q := `INSERT INTO opening_hours (id, place_id, day, open_time, close_time)
				  VALUES ($1, $2, $3, $4, $5)`
		_, err := tx.ExecContext(ctx, q, hour.ID, hour.PlaceId, hour.Day, hour.OpenTime, hour.CloseTime)
		if err != nil {
			return utils.ParsePQError(err)
		}
	}

	return nil
}
