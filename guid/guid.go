package guid

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func New() string {
	return strings.ReplaceAll(uuid.NewV4().String(), "-", "")
}
func NewWith() string {
	return uuid.NewV4().String()
}
