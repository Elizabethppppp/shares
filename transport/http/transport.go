package http

import (
	"CA_shares/core/CA_shares"
	"CA_shares/middleware"
	"CA_shares/storage"
	"crypto/rsa"
	"net/http"
)

type Transport struct {
	s         *CA_shares.SharesService
	publicKey *rsa.PublicKey
	storage   *storage.ImageManager
}

func NewTransport(s *CA_shares.SharesService, publicKey *rsa.PublicKey, storage *storage.ImageManager) *Transport {
	return &Transport{s: s, publicKey: publicKey, storage: storage}
}

func (t *Transport) Handler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /api/v1/news", middleware.LoggerMiddleware(http.HandlerFunc(t.ListNews)))
	mux.Handle("GET /api/v1/news/{id}", middleware.LoggerMiddleware(http.HandlerFunc(t.GetNewsById)))

	adminAuth := middleware.AuthMiddleware(t.publicKey, "User_adm")

	mux.Handle("GET /api/v1/admin/news", middleware.LoggerMiddleware(adminAuth(http.HandlerFunc(t.GetAllNewsAdmin))))
	mux.Handle("GET /api/v1/admin/news/{id}", middleware.LoggerMiddleware(adminAuth(http.HandlerFunc(t.NewsByIdAdminRole))))
	mux.Handle("PATCH /api/v1/admin/news/{id}/status", middleware.LoggerMiddleware(adminAuth(http.HandlerFunc(t.NewsUpdateStatus))))
	mux.Handle("DELETE /api/v1/admin/news/{id}", middleware.LoggerMiddleware(adminAuth(http.HandlerFunc(t.NewsDelete))))
	mux.Handle("POST /api/v1/admin/news", middleware.LoggerMiddleware(adminAuth(http.HandlerFunc(t.CreateNews))))
	mux.Handle("PUT /api/v1/admin/news/{id}", middleware.LoggerMiddleware(adminAuth(http.HandlerFunc(t.UpdateNews))))

	return mux
}
