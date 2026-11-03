package model

import "time"

type News struct {
	Id          int32
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
	Id          int32     `json:"id"`
	Title       string    `json:"title"`
	PublishedAt time.Time `json:"publishedAt"`
	PreviewImg  string    `json:"previewImg"`
}
