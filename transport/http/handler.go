package http

import (
	"CA_shares/errorResponse"
	"CA_shares/model"
	"CA_shares/serviceErrors"
	"CA_shares/validation"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

	var isActive *bool
	if s := q.Get("isActive"); s != "" {
		v, err := strconv.ParseBool(s)
		if err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrBadRequest)
			return
		}
		isActive = &v
	}

	query := model.ListQuery{
		Limit:    limit,
		Offset:   offset,
		Title:    q.Get("title"),
		IsActive: isActive,
		SortBy:   q.Get("sortBy"),
		SortDir:  q.Get("sortDir"),
	}

	resp, err := t.s.GetAllNews(r.Context(), query)
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

	if err := t.s.DeleteNotification(r.Context(), id); err != nil {
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

	input := model.NewsCreateInput{
		Title:       r.FormValue("title"),
		Content:     r.FormValue("content"),
		PublishedAt: r.FormValue("publishedAt"),
	}

	if file, header, err := r.FormFile("previewImg"); err == nil {
		defer file.Close()
		input.PreviewImg = &model.FileInput{Filename: header.Filename, Size: header.Size, Reader: file}
	} else if !errors.Is(err, http.ErrMissingFile) {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}
	if file, header, err := r.FormFile("img"); err == nil {
		defer file.Close()
		input.Img = &model.FileInput{Filename: header.Filename, Size: header.Size, Reader: file}
	} else if !errors.Is(err, http.ErrMissingFile) {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
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

func (t *Transport) UpdateNews(w http.ResponseWriter, r *http.Request) {
	id, err := validation.ParseId(r)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	const maxBody = 16 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)

	if err := r.ParseMultipartForm(maxBody); err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsUpdateFailed)
		return
	}

	input := model.NewsUpdateInput{
		Id:          id,
		Title:       r.FormValue("title"),
		Content:     r.FormValue("content"),
		PublishedAt: r.FormValue("publishedAt"),
	}

	if file, header, err := r.FormFile("previewImg"); err == nil {
		defer file.Close()
		input.PreviewImg = &model.FileInput{Filename: header.Filename, Size: header.Size, Reader: file}
	} else if !errors.Is(err, http.ErrMissingFile) {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}
	if file, header, err := r.FormFile("img"); err == nil {
		defer file.Close()
		input.Img = &model.FileInput{Filename: header.Filename, Size: header.Size, Reader: file}
	} else if !errors.Is(err, http.ErrMissingFile) {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}

	news, err := t.s.ToEditNotification(r.Context(), input)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}
