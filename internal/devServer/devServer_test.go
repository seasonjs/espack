package devServer

import (
	"mime"
	"path/filepath"
	"testing"
)

func TestCtx_Add(t *testing.T) {
	ct := mime.TypeByExtension(filepath.Ext("1.html"))
	t.Log(ct)
}
