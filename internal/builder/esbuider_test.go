package builder

import (
	"fmt"
	"github.com/seasonjs/espack/internal/config"
	"github.com/seasonjs/espack/internal/utils"
	"path/filepath"
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
