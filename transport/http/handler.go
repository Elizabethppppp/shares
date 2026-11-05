package http

import (
	"CA_shares/errorResponse"
	"CA_shares/serviceErrors"
	"encoding/json"
	"net/http"
	"strconv"
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
	idText := r.PathValue("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInvalidId)
		return
	}

	news,err:=t.rs.ListById(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}
