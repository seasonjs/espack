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
	//CompressionLevel int

	// Disables parallel gzip.
	//SingleThreaded bool

	readerWrapFn func(io.Reader) (io.Reader, error)
	writerWrapFn func(io.Writer) (io.Writer, error)
}

func (tgz *TarGz) UnTarGz(source, destination string) error {
	tgz.wrapReader()
	return tgz.Tar.UnTar(source, destination)
}

func (tgz *TarGz) IOUnTarGz(reader io.ReadCloser, destination string) error {
	tgz.wrapReader()
	return tgz.Tar.IOUnTar(reader, destination)
}

// 舍弃并发方案
func (tgz *TarGz) wrapReader() {
	var gzr io.ReadCloser

	tgz.Tar.readerWrapFn = func(r io.Reader) (io.Reader, error) {
		var err error
		gzr, err = gzip.NewReader(r)
		return gzr, err
	}
	tgz.Tar.cleanupWrapFn = func() {
		gzr.Close()
	}
}

// DefaultTarGz 针对npm下载的默认设置，一个有问题需要调整下载参数
func DefaultTarGz() *TarGz {
	return &TarGz{
		//CompressionLevel: gzip.DefaultCompression,
		Tar: &Tar{
			MkdirAll:          true,
			OverwriteExisting: true,
			StripComponents:   1,
		},
	}
}
