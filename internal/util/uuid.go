package util

import (
	"github.com/google/uuid"
	"strings"
)

func GetSecondPathParam(path string) string {
	parts := strings.Split(path, "/")
	return parts[2]
}

func GetUuidFromPathSecondPosition(path string) uuid.UUID {
	return getUuidFromPathIndex(path, 2)
}

func getUuidFromPathIndex(path string, index int) uuid.UUID {
	parts := strings.Split(path, "/")
	return uuid.MustParse(parts[index])
}
