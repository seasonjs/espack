package mod

// 截取修改自：
//https://github.com/mholt/archiver/blob/master/targz.go
import (
	"compress/gzip"
	"io"
)

// TarGz facilitates gzip compression
// (RFC 1952) of tarball archives.
type TarGz struct {
	*Tar

	// The compression level to use, as described
	// in the compress/gzip package.
	CompressionLevel int

	// Disables parallel gzip.
	//SingleThreaded bool
}

func (tgz *TarGz) UnTarGz(source, destination string) error {
	tgz.wrapReader()
	return tgz.Tar.UnTar(source, destination)
}

// 舍弃并发方案
func (tgz *TarGz) wrapReader() {
	var gzr io.ReadCloser
	tgz.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		var err error
		//if tgz.SingleThreaded {
		gzr, err = gzip.NewReader(r)
		//} else {
		//	gzr, err = pgzip.NewReader(r)
		//}
		return gzr, err
	}
	tgz.Tar.cleanupWrapFn = func() {
		gzr.Close()
	}
}

func NewTarGz() *TarGz {
	return &TarGz{
		CompressionLevel: gzip.DefaultCompression,
		Tar:              NewTar(),
	}
}
