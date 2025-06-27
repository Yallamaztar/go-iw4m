package iw4m

import (
	"github.com/Yallamaztar/go-iw4m/wrapper"
)

type Utils struct {
	Wrapper *wrapper.IW4MWrapper
}

func NewUtils(w *wrapper.IW4MWrapper) *Utils {
	return &Utils{Wrapper: w}
}

// func (u *Utils) DoesRoleExists(role string) string {
// 	server := NewServer(u.Wrapper)
// }

// func (u *Utils) RolePosition(role string)
