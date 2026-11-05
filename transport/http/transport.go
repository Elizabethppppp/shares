package http

import (
	"CA_shares/core/CA_shares"
	"CA_shares/middleware"
	"net/http"
)

type Transport struct {
	rs *CA_shares.ReduceService
}

func NewTransport(rs *CA_shares.ReduceService) *Transport {
	return &Transport{rs: rs}
}

func (t *Transport) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/news", middleware.LoggerMiddleware(http.HandlerFunc(t.ListNews)))
	mux.Handle("GET /api/v1/news/{id}", middleware.LoggerMiddleware(http.HandlerFunc(t.GetNewsById)))
	mux.Handle("GET /api/v1/admin/news/{id}", middleware.LoggerMiddleware(http.HandlerFunc(t.NewsByIdAdminRole)))

	return mux
}
