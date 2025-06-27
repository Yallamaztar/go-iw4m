package iw4m

import (
	"net/http"

	"github.com/Yallamaztar/go-iw4m/wrapper"
)

// Constructor to create IW4MWrapper instance
func NewWrapper(baseUrl string, serverID string, cookie string) *wrapper.IW4MWrapper {
	return &wrapper.IW4MWrapper{
		BaseURL:  baseUrl,
		ServerID: serverID,
		Cookie:   cookie,
		Client:   &http.Client{},
	}
}
