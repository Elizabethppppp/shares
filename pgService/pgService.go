package pgService

import (
	"CA_shares/model"
	"CA_shares/serviceErrors"
	"context"
	"database/sql"
	"errors"
	"time"
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
	const query = `SELECT id,title,published_at, 
       				COALESCE(preview_img , '') AS preview_img 
				    FROM news 
					WHERE is_active = TRUE AND published_at <= now() 
					ORDER BY published_at DESC 
                    LIMIT 10 OFFSET 0`
	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		return nil, serviceErrors.ErrInternal
	}
	defer rows.Close()
	var news []model.NewsListPublic
	var publishedAt time.Time
	for rows.Next() {
		var notification model.NewsListPublic
		err := rows.Scan(
			&notification.Id,
			&notification.Title,
			&publishedAt,
			&notification.PreviewImg)
		if err != nil {
			return nil, serviceErrors.ErrInternal
		}
		notification.PublishedAt = publishedAt.Format("2006-01-02")
		news = append(news, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, serviceErrors.ErrInternal
	}
	return news, nil
}

func (pg *PgService) GetById(ctx context.Context, id int64) (model.NewsDetailPublic, error) {
	const query = `SELECT id,title,content, published_at, 
       				COALESCE(img, '') AS img
					FROM news 
					WHERE id = $1
					AND is_active = true 
					AND published_at <= CURRENT_DATE;`

	row := pg.db.QueryRowContext(ctx, query, id)
	var news model.NewsDetailPublic
	var publishedAt time.Time
	err := row.Scan(
		&news.Id,
		&news.Title,
		&news.Content,
		&publishedAt,
		&news.Img)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.NewsDetailPublic{}, serviceErrors.ErrNotFound
		}
		return model.NewsDetailPublic{}, serviceErrors.ErrInternal
	}
	news.PublishedAt = publishedAt.Format("2006-01-02")
	return news, nil
}

func (pg *PgService) GetByIdAdminRole(ctx context.Context, id int64) (model.NewsDetailAdmin, error) {
	const query = `SELECT id,title,content, is_active, published_at, created_at, updated_at, 
       				COALESCE(preview_img, '') AS preview_img, 
       				COALESCE(img, '') AS img
					FROM news 
					WHERE id = $1;`

	row := pg.db.QueryRowContext(ctx, query, id)
	var news model.NewsDetailAdmin
	var publishedAt time.Time
	var createdAt time.Time
	var updatedAt sql.NullTime
	err := row.Scan(
		&news.Id,
		&news.Title,
		&news.Content,
		&news.IsActive,
		&publishedAt,
		&createdAt,
		&updatedAt,
		&news.PreviewImg,
		&news.Img)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.NewsDetailAdmin{}, serviceErrors.ErrNotFound
		}
		return model.NewsDetailAdmin{}, serviceErrors.ErrInternal
	}
	news.PublishedAt = publishedAt.Format("2006-01-02")
	news.CreatedAt = createdAt
	if updatedAt.Valid {
		t := updatedAt.Time
		news.UpdatedAt = &t
	}
	return news, nil
}
