package model

import "time"

type News struct {
	Id          int64
	Title       string
	Content     string
	IsActive    bool
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	PreviewImg  string
	Img         string
}

type NewsListPublic struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
	PreviewImg  string `json:"previewImg"`
}

type NewsListPublicResponse struct {
	NextCursor string           `json:"nextCursor"`
	List       []NewsListPublic `json:"list"`
}

type NewsDetailPublic struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	PublishedAt string `json:"publishedAt"`
	Img         string `json:"img"`
}

type NewsDetailAdmin struct {
	Id          int64   `json:"id"`
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	IsActive    bool    `json:"isActive"`
	PublishedAt string  `json:"publishedAt"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   *string `json:"updatedAt"`
	PreviewImg  string  `json:"previewImg"`
	Img         string  `json:"img"`
}
