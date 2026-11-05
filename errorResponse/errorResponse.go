package errorResponse

import (
	"CA_shares/serviceErrors"
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorCode string

const (
	CodeBadRequest ErrorCode = "400"
	CodeNotFound   ErrorCode = "404"
	CodeInternal   ErrorCode = "500"
	CodeForbidden  ErrorCode = "403"
)

type UserMessage struct {
	LangRu string `json:"lang_ru"`
	LangEn string `json:"lang_en"`
}

type ErrorResponse struct {
	Code         ErrorCode   `json:"code"`
	InternalCode string      `json:"internalCode"`
	DevMessage   string      `json:"devMessage"`
	UserMessage  UserMessage `json:"userMessage"`
}

type errorResponse struct {
	StatusCode   int
	Code         ErrorCode
	InternalCode string
	DevMessage   string
	UserMessage  UserMessage
}

var errorResponseMessage = map[error]errorResponse{
	serviceErrors.ErrNotFound: {
		StatusCode:   http.StatusNotFound,
		Code:         CodeNotFound,
		InternalCode: "",
		DevMessage:   "",
		UserMessage: UserMessage{
			LangRu: "Уведомление не найдено.",
			LangEn: "Notification not found.",
		},
	},
	serviceErrors.ErrBadRequest: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "INVALID_REQUEST_DATA",
		DevMessage:   "Invalid Request Data",
		UserMessage: UserMessage{
			LangRu: "Некорректные данные запроса",
			LangEn: "Invalid Request Data",
		},
	},
	serviceErrors.ErrInternal: {
		StatusCode:   http.StatusInternalServerError,
		Code:         CodeInternal,
		InternalCode: "INTERNAL_ERROR",
		DevMessage:   "Internal Error",
		UserMessage: UserMessage{
			LangRu: "Внутренняя ошибка сервера",
			LangEn: "Internal Error",
		},
	},
	serviceErrors.ErrInvalidId: {
		StatusCode:   http.StatusBadRequest,
		Code:         CodeBadRequest,
		InternalCode: "INVALID_ID",
		DevMessage:   "Invalid Id format",
		UserMessage: UserMessage{
			LangEn: "Invalid Id format",
			LangRu: "Некорректный формат Id",
		},
	},
	serviceErrors.ErrForbidden: {
		StatusCode:   http.StatusForbidden,
		Code:         CodeForbidden,
		InternalCode: "FORBIDDEN",
		DevMessage:   "Forbidden",
		UserMessage: UserMessage{
			LangEn: "Access Denied",
			LangRu: "Доступ запрещён",
		},
	},
}

func ErrorResponseJSON(w http.ResponseWriter, errs error) {
	info, ok := errorResponseMessage[serviceErrors.ErrInternal]
	for target, mapped := range errorResponseMessage {
		if errors.Is(errs, target) {
			info = mapped
			ok = true
			break
		}
	}
	if !ok {
		info = errorResponseMessage[serviceErrors.ErrInternal]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(info.StatusCode)

	resp := ErrorResponse{
		Code:         info.Code,
		InternalCode: info.InternalCode,
		DevMessage:   info.DevMessage,
		UserMessage:  info.UserMessage,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		_, _ = w.Write([]byte(`{"code":"500","internalCode":"MARSHAL_ERROR","devMessage":"failed to marshal error","userMessage":{"lang_ru":"Внутренняя ошибка сервера","lang_en":"Internal server error"}}`))
		return
	}

	_, _ = w.Write(body)
}
