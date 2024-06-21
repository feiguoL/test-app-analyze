package utils

import (
	"time"

	"github.com/lytdev/go-mykit/guid"
)

const gotime = "2006-01-02T15:04:05"

func GetCurrentTimeStr() string {
	now := time.Now()
	formatterTime := now.Format(gotime)
	return formatterTime
}

func GetFastUuid() string {
	return guid.FastUuid()
}
