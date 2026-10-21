package pkg

import (
	"github.com/google/uuid"
)

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}
