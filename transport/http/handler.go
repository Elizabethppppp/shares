package http

import (
	"CA_shares/errorResponse"
	"CA_shares/validation"
	"encoding/json"
	"net/http"
)

func (t *Transport) ListNews(w http.ResponseWriter, r *http.Request) {
	resp, err := t.rs.List(r.Context())
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (t *Transport) GetNewsById(w http.ResponseWriter, r *http.Request) {
	id, err := validation.ParseId(r)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	news, err := t.rs.NotificationById(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}

func (t *Transport) NewsByIdAdminRole(w http.ResponseWriter, r *http.Request) {
	id, err := validation.ParseId(r)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	news, err := t.rs.NotificationByIdAdminRole(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}

func (t *Transport) NewsUpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := validation.ParseId(r)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	news, err := t.rs.UpdateStatus(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}

func (t *Transport) NewsDelete(w http.ResponseWriter, r *http.Request) {
	id, err := validation.ParseId(r)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	err = t.rs.DeleteNotification(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
