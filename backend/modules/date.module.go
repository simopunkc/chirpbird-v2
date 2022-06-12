package modules

import (
	"time"
)

func GenerateIsoDate() string {
	return time.Now().Format(time.RFC3339Nano)
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
