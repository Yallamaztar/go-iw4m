package http

import (
	"io"
	"net/http"
)

// Wrapper For Making Request
func (w *IW4MWrapper) DoRequest(path string) string {
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Cookie", w.Cookie)

	r, err := w.Client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer r.Body.Close()

	body, _ := io.ReadAll(r.Body)
	return string(body)
}
