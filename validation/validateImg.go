package validation

import (
	"CA_shares/serviceErrors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func ValidateImgFile(header *multipart.FileHeader, maxSize int64) error {
	if header.Size > maxSize {
		return serviceErrors.ErrFileTooLarge
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	switch ext {
	case ".jpg", ".jpeg", ".png":
	default:
		return fmt.Errorf("invalid file type: %s", ext)
	}

	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return err
	}

	mimeType := http.DetectContentType(buffer[:n])
	switch mimeType {
	case "image/jpeg", "image/jpg", "image/png":
	default:
		return fmt.Errorf("invalid file type: %s", mimeType)
	}
	return nil
}
