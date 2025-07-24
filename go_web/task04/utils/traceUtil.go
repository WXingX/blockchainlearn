package utils

import "github.com/google/uuid"

func GenTraceID() string {
	return uuid.NewString()
}
