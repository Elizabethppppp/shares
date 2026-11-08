package http

import (
	"CA_shares/core/CA_shares"
	"CA_shares/middleware"
	"net/http"
)

type Transport struct {
	s *CA_shares.SharesService
}

func NewTransport(s *CA_shares.SharesService) *Transport {
	return &Transport{s: s}
}

func (t *Transport) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/news", middleware.LoggerMiddleware(http.HandlerFunc(t.ListNews)))
	mux.Handle("GET /api/v1/news/{id}", middleware.LoggerMiddleware(http.HandlerFunc(t.GetNewsById)))
	mux.Handle("GET /api/v1/admin/news", middleware.LoggerMiddleware(http.HandlerFunc(t.GetAllNewsAdmin)))
	mux.Handle("GET /api/v1/admin/news/{id}", middleware.LoggerMiddleware(http.HandlerFunc(t.NewsByIdAdminRole)))
	mux.Handle("PATCH /api/v1/admin/news/{id}/status", middleware.LoggerMiddleware(http.HandlerFunc(t.NewsUpdateStatus)))
	mux.Handle("DELETE /api/v1/admin/news/{id}", middleware.LoggerMiddleware(http.HandlerFunc(t.NewsDelete)))


	return mux
}
