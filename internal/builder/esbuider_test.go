package builder

import (
	"fmt"
	"path/filepath"
	"seasonjs/espack/internal/config"
	"seasonjs/espack/internal/utils"
	"testing"
)

func TestEsbuild_GetOptions(t *testing.T) {
	path := utils.FS.GetCurrentPath()
	path = filepath.Join(path, "../case/espack.config.json")
	cf := config.
		NewConfigPoints().
		ReadFile(path).
		ReadConfig()
	opt := NewEsBuilder(*cf).GetOptions().EsbuildStarter()
	fmt.Println(opt)
}
