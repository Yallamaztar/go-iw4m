package main

import (
	"net/http"
	"os"

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

func main() {
	wrapper := NewWrapper(
		os.Getenv("IW4M_URL"),
		os.Getenv("IW4M_ID"),
		os.Getenv("IW4M_HEADER"),
	)

	server := iw4m.NewServer(wrapper)
}
