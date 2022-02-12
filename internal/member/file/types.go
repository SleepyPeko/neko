package file

import (
	"demodesk/neko/internal/types"
)

type MemberEntry struct {
	Password string              `json:"password"`
	Profile  types.MemberProfile `json:"profile"`
}

type Config struct {
	Path string
}