package utils

import (
	"github.com/google/uuid"
	"time"
)

const gotime = "2006-01-02T15:04:05"

func GetCurrentTimeStr() string {
	now := time.Now()
	formatterTime := now.Format(gotime)
	return formatterTime
}

func GetFastUuid() string {
	return uuid.New().String()
}
