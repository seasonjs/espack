package utils

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

// 无实际意义，仅作为区分使用
type env struct {
}

func newEnv() {
	path, _ := filepath.Abs("./internal/utils/env/.env")
	e := godotenv.Load(path)
	if e != nil {
		Err.LogAndExit(e)
	}
}
func (e env) Dev() bool {
	return os.Getenv("ES_PACK_ENV") == "dev"
}

func (e env) Prod() bool {
	return os.Getenv("ES_PACK_ENV") == "prod"
}
