package main

import (
	"encoding/json"
	"fmt"
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
