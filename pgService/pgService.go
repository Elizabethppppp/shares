package pgService

import (
	"CA_shares/model"
	"CA_shares/serviceErrors"
	"context"
	"database/sql"
)

type PgService struct {
	db *sql.DB
}

func NewPgService(db *sql.DB) *PgService {
	return &PgService{
		db: db,
	}
}

func (pg *PgService) GetAll(ctx context.Context) ([]model.NewsListPublic, error) {
	const query = `SELECT id,title,published_at, preview_img FROM news WHERE is_active = TRUE AND published_at <= now() ORDER BY published_at DESC LIMIT 10 OFFSET 0`
	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		return nil, serviceErrors.ErrInternal
	}
	defer rows.Close()
	var news []model.NewsListPublic
	for rows.Next() {
		var notification model.NewsListPublic
		err := rows.Scan(
			&notification.Id,
			&notification.Title,
			&notification.PublishedAt,
			&notification.PreviewImg)
		if err != nil {
			return nil, serviceErrors.ErrInternal
		}
		news = append(news, notification)
	}
	if err := rows.Err(); err != nil {
		return nil, serviceErrors.ErrInternal
	}
	return news, nil
}
