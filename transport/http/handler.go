package http

import (
	"CA_shares/errorResponse"
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
