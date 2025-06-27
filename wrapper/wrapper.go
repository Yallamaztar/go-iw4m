package wrapper

import "net/http"

type IW4MWrapper struct {
	BaseURL  string
	ServerID int
	Cookie   string
	Client   *http.Client
}
