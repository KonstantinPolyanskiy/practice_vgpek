package apiutils

import (
	"fmt"
	"net/http"
)

// SetDownloadHeaders устанавливает нужные заголовки при отдаче файла с сервера
func SetDownloadHeaders(w http.ResponseWriter, name, len string) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, name))
	w.Header().Set("Content-Length", len)
	w.Header().Set("Cache-Control", "private")
	w.Header().Set("Pragma", "private")
}
