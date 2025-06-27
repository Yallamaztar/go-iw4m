package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type IW4MWrapper struct {
	BaseURL  string
	ServerID int
	Cookie   string
	Client   *http.Client
}

// Return New Instance Of The Wrapper
func NewWrapper(baseUrl string, serverID int, cookie string) *IW4MWrapper {
	return &IW4MWrapper{
		BaseURL:  baseUrl,
		ServerID: serverID,
		Cookie:   cookie,
		Client:   &http.Client{},
	}
}

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

func main() {
	serverIDStr := os.Getenv("IW4M_ID")
	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		fmt.Println("Invalid IW4M_ID:", err)
		return
	}

	wrapper := NewWrapper(
		os.Getenv("IW4M_URL"),
		serverID,
		os.Getenv("IW4M_HEADER"),
	)
	rules := wrapper.Server().Rules()
	b, _ := json.MarshalIndent(rules, "", "  ")
	fmt.Println(string(b))
}
