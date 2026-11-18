package validation

import (
	"CA_shares/serviceErrors"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func ValidateImgFile(fileName string, size int64, r io.Reader, maxSize int64) (io.Reader, error) {
	if size > maxSize {
		return nil, serviceErrors.ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg", ".png":
	default:
		return nil, fmt.Errorf("invalid file type: %s", ext)
	}

	buffer := make([]byte, 512)
	n, err := r.Read(buffer)
	if err != nil && n == 0 {
		return nil, err
	}

	mimeType := http.DetectContentType(buffer[:n])
	switch mimeType {
	case "image/jpeg", "image/png":
	default:
		return nil, fmt.Errorf("invalid file type: %s", mimeType)
	}
	return io.MultiReader(bytes.NewReader(buffer[:n]), r), nil
}
