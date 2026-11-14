package http

import (
	"CA_shares/errorResponse"
	"CA_shares/filtres"
	"CA_shares/logger"
	"CA_shares/model"
	"CA_shares/serviceErrors"
	"CA_shares/validation"
	"encoding/json"
	"errors"
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
	previewImg, img, err := t.s.DeleteNotification(r.Context(), id)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	if err := t.storage.Delete(previewImg); err != nil {
		logger.Error("delete preview file", "key", previewImg, "error", err.Error())
	}

	if err := t.storage.Delete(img); err != nil {
		logger.Error("delete img file", "key", img, "error", err.Error())
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
	if content == "" || len([]rune(content)) > 10000 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}

	pubAtStr := r.FormValue("publishedAt")
	publishedAt, err := time.Parse("02-01-2006", pubAtStr)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
		return
	}

	var previewKey, imgKey string

	if file, header, err := r.FormFile("previewImg"); err == nil {
		defer file.Close()
		if err := validation.ValidateImgFile(header, 5<<20); err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
			return
		}
		key, err := t.storage.Save("previewImg", header.Filename, file)
		if err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInternal)
			return
		}
		previewKey = key
	}

	if file, header, err := r.FormFile("img"); err == nil {
		defer file.Close()
		if err := validation.ValidateImgFile(header, 10<<20); err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsValidationFailed)
			return
		}
		key, err := t.storage.Save("img", header.Filename, file)
		if err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInternal)
			return
		}
		imgKey = key
	}

	input := model.NewsPostNew{
		Title:       title,
		Content:     content,
		PublishedAt: publishedAt,
		PreviewImg:  previewKey,
		Img:         imgKey,
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

	title := r.FormValue("title")
	if title == "" || len([]rune(title)) > 50 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsUpdateFailed)
		return
	}

	content := r.FormValue("content")
	if content == "" || len([]rune(content)) > 10000 {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsUpdateFailed)
		return
	}

	pubAtStr := r.FormValue("publishedAt")
	publishedAt, err := time.Parse("02-01-2006", pubAtStr)
	if err != nil {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsUpdateFailed)
		return
	}

	var newPreviewImg, newImg *string

	if file, header, err := r.FormFile("previewImg"); err == nil {
		defer file.Close()
		if err := validation.ValidateImgFile(header, 5<<20); err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsUpdateFailed)
			return
		}
		key, err := t.storage.Save("previewImg", header.Filename, file)
		if err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInternal)
			return
		}
		newPreviewImg = &key
	} else if !errors.Is(err, http.ErrMissingFile) {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsUpdateFailed)
		return
	}

	if file, header, err := r.FormFile("img"); err == nil {
		defer file.Close()
		if err := validation.ValidateImgFile(header, 10<<20); err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsUpdateFailed)
			return
		}
		key, err := t.storage.Save("img", header.Filename, file)
		if err != nil {
			errorResponse.ErrorResponseJSON(w, serviceErrors.ErrInternal)
			return
		}
		newImg = &key
	} else if !errors.Is(err, http.ErrMissingFile) {
		errorResponse.ErrorResponseJSON(w, serviceErrors.ErrNewsUpdateFailed)
		return
	}

	input := model.NewsEditNotification{
		Id:          id,
		Title:       title,
		Content:     content,
		PublishedAt: publishedAt,
		PreviewImg:  newPreviewImg,
		Img:         newImg,
	}

	news, oldPreviewImg, oldImg, err := t.s.ToEditNotification(r.Context(), input)
	if err != nil {
		if newPreviewImg != nil {
			_ = t.storage.Delete(oldPreviewImg)
		}
		if newImg != nil {
			_ = t.storage.Delete(oldImg)
		}
		errorResponse.ErrorResponseJSON(w, err)
		return
	}

	if err := t.storage.Delete(oldPreviewImg); err != nil {
		logger.Error("delete old preview", "key", oldPreviewImg, "error", err.Error())
	}
	if err := t.storage.Delete(oldImg); err != nil {
		logger.Error("delete old img", "key", oldImg, "error", err.Error())
	}

	w.Header().Set("Content-Type", "multipart/form-data")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(news)
}
