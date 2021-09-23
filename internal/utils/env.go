package utils

import (
	"os"
)

// 无实际意义，仅作为区分使用
type env struct {
}

func (e env) Dev() bool {
	return os.Getenv("ES_PACK_DEV_ENV") == "dev"
}

func (e env) Prod() bool {
	return os.Getenv("ES_PACK_DEV_ENV") == "prod"
}
