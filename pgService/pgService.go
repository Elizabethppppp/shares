package pgService

import (
	"CA_shares/filtres"
	"CA_shares/logger"
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

func (pg *PgService) GetAllList(ctx context.Context, f filtres.ListFilter) ([]model.NewsList, int64, error) {
	orderBy := f.SortBy
	var orderDir string
	switch f.SortDir {
	case "asc":
		orderDir = "ASC"
	case "desc":
		orderDir = "DESC"
	default:
		orderDir = "DESC"
	}
	query := `SELECT id, title, is_active, published_at, created_at,
               COALESCE(preview_img, '') AS preview_img,
               COALESCE(img, '') AS img
     	   	   FROM news
        	   WHERE ($1::text IS NULL OR title ILIKE '%' || $1 || '%')
               AND ($2::boolean IS NULL OR is_active = $2)
               ORDER BY ` + orderBy + ` ` + orderDir + `, id ` + orderDir + `
               LIMIT $3 OFFSET $4`
	rows, err := pg.db.QueryContext(ctx, query, f.Title, f.IsActive, f.Limit, f.Offset)
	if err != nil {
		return nil, 0, serviceErrors.ErrInternal
	}
	defer rows.Close()

	var items []model.NewsList
	for rows.Next() {
		var item model.NewsList
		var publishedAt time.Time
		var createdAt time.Time
		err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.IsActive,
			&publishedAt,
			&createdAt,
			&item.PreviewImg,
			&item.Img)
		if err != nil {
			return nil, 0, serviceErrors.ErrInternal
		}
		item.PublishedAt = publishedAt.Format("2006-01-02")
		item.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, serviceErrors.ErrInternal
	}

	var count int64
	const countQuery = `SELECT COUNT(*) FROM news
						WHERE ($1::text IS NULL OR title ILIKE '%' || $1 || '%')
          				AND ($2::boolean IS NULL OR is_active = $2)`
	if err := pg.db.QueryRowContext(ctx, countQuery, f.Title, f.IsActive).Scan(&count); err != nil {
		return nil, 0, serviceErrors.ErrInternal
	}
	return items, count, nil
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
	news.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	if updatedAt.Valid {
		t := updatedAt.Time.Format("2006-01-02 15:04:05")
		news.UpdatedAt = &t
	}
	return news, nil
}

func (pg *PgService) UpdateActive(ctx context.Context, id int64) (model.NewsUpdateActive, error) {
	var publishedAt time.Time
	err := pg.db.QueryRowContext(ctx,
		`SELECT published_at FROM news WHERE id = $1`, id).Scan(&publishedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return model.NewsUpdateActive{}, serviceErrors.ErrNotFound
	}
	if err != nil {
		return model.NewsUpdateActive{}, serviceErrors.ErrInternal
	}
	if publishedAt.After(time.Now()) {
		return model.NewsUpdateActive{}, serviceErrors.ErrPublicationNotReached
	}

	const query = `UPDATE news
        			SET is_active = NOT is_active, updated_at = NOW()
     			    WHERE id = $1
        			RETURNING id, title, content, is_active, published_at, created_at, updated_at,
                 	COALESCE(preview_img, '') AS preview_img,
                  	COALESCE(img, '') AS img`

	var news model.NewsUpdateActive
	var pubAt time.Time
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = pg.db.QueryRowContext(ctx, query, id).Scan(
		&news.Id,
		&news.Title,
		&news.Content,
		&news.IsActive,
		&pubAt,
		&createdAt,
		&updatedAt,
		&news.PreviewImg,
		&news.Img,
	)
	if err != nil {
		return model.NewsUpdateActive{}, serviceErrors.ErrInternal
	}

	news.PublishedAt = pubAt.Format("2006-01-02")
	news.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	if updatedAt.Valid {
		t := updatedAt.Time.Format("2006-01-02 15:04:05")
		news.UpdatedAt = &t
	}
	return news, nil
}

func (pg *PgService) Delete(ctx context.Context, id int64) (string, string, error) {
	const query = `DELETE FROM news WHERE id = $1
					RETURNING COALESCE(preview_img, ''), COALESCE(img, '')`
	var previewImg, img string

	err := pg.db.QueryRowContext(ctx, query, id).Scan(&previewImg, &img)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", serviceErrors.ErrNotFound
	}
	if err != nil {
		return "", "", serviceErrors.ErrInternal
	}

	return previewImg, img, nil
}

func (pg *PgService) Create(ctx context.Context, p model.NewsPostNew) (model.NewsList, error) {
	const query = `INSERT INTO news (title, content, published_at, preview_img, img)
					VALUES ($1, $2, $3, $4, $5)
					RETURNING id, title, is_active, published_at, created_at,
          			COALESCE(preview_img, '') AS preview_img,
          			COALESCE(img, '') AS img`

	var news model.NewsList
	var publishedAt time.Time
	var createdAt time.Time

	err := pg.db.QueryRowContext(ctx, query, p.Title, p.Content, p.PublishedAt, p.PreviewImg, p.Img).Scan(
		&news.Id,
		&news.Title,
		&news.IsActive,
		&publishedAt,
		&createdAt,
		&news.PreviewImg,
		&news.Img,
	)
	if err != nil {
		return model.NewsList{}, serviceErrors.ErrInternal
	}
	news.PublishedAt = publishedAt.Format("2006-01-02")
	news.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	return news, nil
}

func (pg *PgService) ToEdit(ctx context.Context, e model.NewsEditNotification) (model.NewsList, string, string, error) {

	var oldPreviewImg, oldImg string
	err := pg.db.QueryRowContext(ctx,
		`SELECT COALESCE(preview_img,''), COALESCE(img,'') FROM news WHERE id = $1`, e.Id,
	).Scan(&oldPreviewImg, &oldImg)
	if errors.Is(err, sql.ErrNoRows) {
		return model.NewsList{}, "", "", serviceErrors.ErrNotFound
	}
	if err != nil {
		logger.Error("get old keys", "id", e.Id, "error", err.Error())
		return model.NewsList{}, "", "", serviceErrors.ErrInternal
	}

	const query = `UPDATE news
					SET title = $2, 
					    content = $3, 
					    published_at = $4,
					    is_active = CASE
							WHEN is_active = TRUE AND $4 > CURRENT_DATE THEN FALSE
							ELSE is_active
					END,
					    preview_img = $5, 
						img = $6,
						updated_at=NOW()
					WHERE id = $1
					RETURNING id, title, is_active, published_at, created_at, updated_at,
          			COALESCE(preview_img, '') AS preview_img,
          			COALESCE(img, '') AS img`

	var news model.NewsList
	var pubAt time.Time
	var createdAt time.Time
	var updatedAt sql.NullTime
	err = pg.db.QueryRowContext(ctx, query, e.Id, e.Title, e.Content, e.PublishedAt, e.PreviewImg, e.Img).Scan(
		&news.Id,
		&news.Title,
		&news.IsActive,
		&pubAt,
		&createdAt,
		&updatedAt,
		&news.PreviewImg,
		&news.Img)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.NewsList{}, "", "", serviceErrors.ErrNotFound
		}
		logger.Error("update news", "id", e.Id, "error", err.Error())
		return model.NewsList{}, "", "", serviceErrors.ErrInternal
	}
	news.PublishedAt = pubAt.Format("2006-01-02")
	news.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	if updatedAt.Valid {
		t := updatedAt.Time.Format("2006-01-02 15:04:05")
		_ = t
	}
	return news, oldPreviewImg, oldImg, nil
}
