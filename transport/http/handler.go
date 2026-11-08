package http

import (
	"CA_shares/errorResponse"
	"CA_shares/filtres"
	"CA_shares/model"
	"CA_shares/serviceErrors"
	"CA_shares/validation"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (t *Transport) ListNews(w http.ResponseWriter, r *http.Request) {
	resp, err := t.s.List(r.Context())
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

	news, err := t.s.NotificationById(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}

func (t *Transport) GetAllNewsAdmin(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	offset := 0
	if s := q.Get("offset"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrBadRequest)
			return
		}
		offset = v
	}

	limit := 10
	if s := q.Get("limit"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrBadRequest)
			return
		}
		limit = v
	}

	var title *string
	if s := q.Get("title"); s != "" {
		if len([]rune(s)) > 50 {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrBadRequest)
			return
		}
		title = &s
	}

	var isActive *bool
	if s := q.Get("isActive"); s != "" {
		v, err := strconv.ParseBool(s)
		if err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrBadRequest)
			return
		}
		isActive = &v
	}

	sortBy := "published_at"
	if s := q.Get("sortBy"); s != "" {
		switch s {
		case "publishedAt":
			sortBy = "published_at"
		case "id":
			sortBy = "id"
		default:
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrBadRequest)
			return
		}
	}

	sortDir := "desc"
	if s := q.Get("sortDir"); s != "" {
		switch s {
		case "asc", "desc":
			sortDir = s
		default:
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrBadRequest)
			return
		}
	}

	filter := filtres.ListFilter{
		Offset:   offset,
		Limit:    limit,
		Title:    title,
		IsActive: isActive,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}

	resp, err := t.s.GetAllNews(r.Context(), filter)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

func (t *Transport) NewsByIdAdminRole(w http.ResponseWriter, r *http.Request) {
	id, err := validation.ParseId(r)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	news, err := t.s.NotificationByIdAdminRole(r.Context(), id)
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
	news, err := t.s.UpdateStatus(r.Context(), id)
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
	err = t.s.DeleteNotification(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (t *Transport) CreateNews(w http.ResponseWriter, r *http.Request) {
	const maxBody = 16 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)

	if err := r.ParseMultipartForm(maxBody); err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}

	title := r.FormValue("title")
	if title == "" || len([]rune(title)) > 50 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}

	content := r.FormValue("content")
	if content == "" || len([]rune(content)) > 1000 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}

	pubAtStr := r.FormValue("publishedAt")
	publishedAt, err := time.Parse("02-01-2006", pubAtStr)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}

	input := model.NewsPostNew{
		Title:       title,
		Content:     content,
		PublishedAt: publishedAt,
		PreviewImg:  "",
		Img:         "",
	}

	news, err := t.s.CreateNotification(r.Context(), input)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}
