package services

import (
	"github.com/Yallamaztar/go-iw4m/wrapper"
)

type Utils struct {
	Wrapper *wrapper.IW4MWrapper
}

func NewUtils(w *wrapper.IW4MWrapper) *Utils {
	return &Utils{Wrapper: w}
}

// func DoesRoleExists(role string) bool {
// 	// server := NewServer()
// }
