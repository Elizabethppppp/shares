package validation

import (
	"CA_shares/serviceErrors"
	"net/http"
	"strconv"
)

func ParseId(r *http.Request) (int64, error) {
	idText := r.PathValue("id")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		return 0, serviceErrors.ErrInvalidId
	}
	return id, nil
}
