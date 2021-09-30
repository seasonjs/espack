package mod

//截取修改自:
//https://github.com/mholt/archiver/blob/master/tar.go
import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"seasonjs/espack/internal/logger"
	"strings"
)

//==================== archiver===================

// File provides methods for accessing information about
// or contents of a file within an archive.
type File struct {
	os.FileInfo

	// The original header info; depends on
	// type of archive -- could be nil, too.
	Header interface{}

	// Allow the file contents to be read (and closed)
	io.ReadCloser
}

// ReadFakeCloser is an io.Reader that has
// a no-op close method to satisfy the
// io.ReadCloser interface.
type ReadFakeCloser struct {
	io.Reader
}

// Close implements io.Closer.
func (rfc ReadFakeCloser) Close() error { return nil }

//=========================error===================================

// IllegalPathError is an error returned when an illegal
// path is detected during the archival process.
//
// By default, only the Filename is showed on error, but you might
// also get the absolute value of the invalid path on the AbsolutePath
// field.
type IllegalPathError struct {
	AbsolutePath string
	Filename     string
}

func (err *IllegalPathError) Error() string {
	return fmt.Sprintf("illegal file path: %s", err.Filename)
}

// IsIllegalPathError returns true if the provided error is of
// the type IllegalPathError.
func IsIllegalPathError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "illegal file path: ")
}

// Tar provides facilities for operating TAR archives.
// See http://www.gnu.org/software/tar/manual/html_node/Standard.html.
type Tar struct {
	// Whether to overwrite existing files; if false,
	// an error is returned if the file exists.
	OverwriteExisting bool

	// Whether to make all the directories necessary
	// to create a tar archive in the desired path.
	MkdirAll bool

	// A single top-level folder can be implicitly
	// created by the Archive or Unarchive methods
	// if the files to be added to the archive
	// or the files to be extracted from the archive
	// do not all have a common root. This roughly
	// mimics the behavior of archival tools integrated
	// into OS file browsers which create a subfolder
	// to avoid unexpectedly littering the destination
	// folder with potentially many files, causing a
	// problematic cleanup/organization situation.
	// This feature is available for both creation
	// and extraction of archives, but may be slightly
	// inefficient with lots and lots of files,
	// especially on extraction.
	ImplicitTopLevelFolder bool

	// Strip number of leading paths. This feature is available
	// only during unpacking of the entire archive.
	StripComponents int

	// If true, errors encountered during reading
	// or writing a single file will be logged and
	// the operation will continue on remaining files.
	ContinueOnError bool

	tw *tar.Writer
	tr *tar.Reader

	readerWrapFn  func(io.Reader) (io.Reader, error)
	writerWrapFn  func(io.Writer) (io.Writer, error)
	cleanupWrapFn func()
}

// CheckPath ensures that the filename has not been crafted to perform path traversal attacks
func (*Tar) CheckPath(to, filename string) error {
	to, _ = filepath.Abs(to) //explicit the destination folder to prevent that 'string.HasPrefix' check can be 'bypassed' when no destination folder is supplied in input
	dest := filepath.Join(to, filename)
	//prevent path traversal attacks
	if !strings.HasPrefix(dest, to) {
		return &IllegalPathError{AbsolutePath: dest, Filename: filename}
	}
	return nil
}
func (t *Tar) UnTar(source, destination string) error {
	if !fileExists(destination) && t.MkdirAll {
		err := mkdir(destination, 0755)
		if err != nil {
			return fmt.Errorf("preparing destination: %v", err)
		}
	}

	// if the files in the archive do not all share a common
	// root, then make sure we extract to a single subfolder
	// rather than potentially littering the destination...
	if t.ImplicitTopLevelFolder {
		var err error
		destination, err = t.addTopLevelFolder(source, destination)
		if err != nil {
			return fmt.Errorf("scanning source archive: %v", err)
		}
	}

	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("opening source archive: %v", err)
	}
	defer file.Close()

	err = t.Open(file, 0)
	if err != nil {
		return fmt.Errorf("opening tar archive for reading: %v", err)
	}
	defer t.Close()

	for {
		err := t.untarNext(destination)
		if err == io.EOF {
			break
		}
		if err != nil {
			if t.ContinueOnError || IsIllegalPathError(err) {
				logger.Warn("在读取压缩文件过程中出现错误: %v", err)
				continue
			}
			return fmt.Errorf("reading file in tar archive: %v", err)
		}
	}

	return nil
}
func (t *Tar) IOUnTar(reader io.ReadCloser, destination string) error {
	if !fileExists(destination) && t.MkdirAll {
		err := mkdir(destination, 0755)
		if err != nil {
			return fmt.Errorf("preparing destination: %v", err)
		}
	}
	err := t.Open(reader, 0)
	if err != nil {
		return fmt.Errorf("opening tar archive for reading: %v", err)
	}
	defer t.Close()

	for {
		err := t.untarNext(destination)
		if err == io.EOF {
			break
		}
		if err != nil {
			if t.ContinueOnError || IsIllegalPathError(err) {
				logger.Warn("在读取压缩文件过程中出现错误: %v", err)
				continue
			}
			return fmt.Errorf("reading file in tar archive: %v", err)
		}
	}

	return nil
}

//addTopLevelFolder scans the files contained inside
// the tarball named sourceArchive and returns a modified
// destination if all the files do not share the same
// top-level folder.
func (t *Tar) addTopLevelFolder(sourceArchive, destination string) (string, error) {
	file, err := os.Open(sourceArchive)
	if err != nil {
		return "", fmt.Errorf("opening source archive: %v", err)
	}
	defer file.Close()

	// if the reader is to be wrapped, ensure we do that now
	// or we will not be able to read the archive successfully
	reader := io.Reader(file)
	if t.readerWrapFn != nil {
		reader, err = t.readerWrapFn(reader)
		if err != nil {
			return "", fmt.Errorf("wrapping reader: %v", err)
		}
	}
	if t.cleanupWrapFn != nil {
		defer t.cleanupWrapFn()
	}

	tr := tar.NewReader(reader)

	var files []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("scanning tarball's file listing: %v", err)
		}
		files = append(files, hdr.Name)
	}

	if multipleTopLevels(files) {
		destination = filepath.Join(destination, folderNameFromFileName(sourceArchive))
	}

	return destination, nil
}

func (t *Tar) untarNext(destination string) error {
	f, err := t.Read()
	if err != nil {
		return err // don't wrap error; calling loop must break on io.EOF
	}
	defer f.Close()

	header, ok := f.Header.(*tar.Header)
	if !ok {
		return fmt.Errorf("expected header to be *tar.Header but was %T", f.Header)
	}

	errPath := t.CheckPath(destination, header.Name)
	if errPath != nil {
		return fmt.Errorf("checking path traversal attempt: %v", errPath)
	}

	if t.StripComponents > 0 {
		if strings.Count(header.Name, "/") < t.StripComponents {
			return nil // skip path with fewer components
		}

		for i := 0; i < t.StripComponents; i++ {
			slash := strings.Index(header.Name, "/")
			header.Name = header.Name[slash+1:]
		}
	}
	return t.untarFile(f, destination, header)
}

func (t *Tar) untarFile(f File, destination string, hdr *tar.Header) error {
	to := filepath.Join(destination, hdr.Name)

	// do not overwrite existing files, if configured
	if !f.IsDir() && !t.OverwriteExisting && fileExists(to) {
		return fmt.Errorf("file already exists: %s", to)
	}

	switch hdr.Typeflag {
	case tar.TypeDir:
		return mkdir(to, f.Mode())
	case tar.TypeReg, tar.TypeRegA, tar.TypeChar, tar.TypeBlock, tar.TypeFifo, tar.TypeGNUSparse:
		return writeNewFile(to, f, f.Mode())
	case tar.TypeSymlink:
		return writeNewSymbolicLink(to, hdr.Linkname)
	case tar.TypeLink:
		return writeNewHardLink(to, filepath.Join(destination, hdr.Linkname))
	case tar.TypeXGlobalHeader:
		return nil // ignore the pax global header from git-generated tarballs
	default:
		return fmt.Errorf("%s: unknown type flag: %c", hdr.Name, hdr.Typeflag)
	}
}

// Open opens t for reading an archive from
// in. The size parameter is not used.
func (t *Tar) Open(in io.Reader, size int64) error {
	if t.tr != nil {
		return fmt.Errorf("tar archive is already open for reading")
	}
	// wrapping readers allows us to open compressed tarballs
	if t.readerWrapFn != nil {
		var err error
		in, err = t.readerWrapFn(in)
		if err != nil {
			return fmt.Errorf("wrapping file reader: %v", err)
		}
	}
	t.tr = tar.NewReader(in)
	return nil
}

// Read reads the next file from t, which must have
// already been opened for reading. If there are no
// more files, the error is io.EOF. The File must
// be closed when finished reading from it.
func (t *Tar) Read() (File, error) {
	if t.tr == nil {
		return File{}, fmt.Errorf("tar archive is not open")
	}

	hdr, err := t.tr.Next()
	if err != nil {
		return File{}, err // don't wrap error; preserve io.EOF
	}

	file := File{
		FileInfo:   hdr.FileInfo(),
		Header:     hdr,
		ReadCloser: ReadFakeCloser{t.tr},
	}

	return file, nil
}

// Close closes the tar archive(s) opened by Create and Open.
func (t *Tar) Close() error {
	var err error
	if t.tr != nil {
		t.tr = nil
	}
	if t.tw != nil {
		tw := t.tw
		t.tw = nil
		err = tw.Close()
	}
	// make sure cleanup of "Reader/Writer wrapper"
	// (say that ten times fast) happens AFTER the
	// underlying stream is closed
	if t.cleanupWrapFn != nil {
		t.cleanupWrapFn()
	}
	return err
}

// NewTar returns a new, default instance ready to be customized and used.
func NewTar() *Tar {
	return &Tar{
		MkdirAll:          true,
		OverwriteExisting: true,
	}
}

//// DeCompress 解压 tar.gz
//func DeCompress(input, output string) error {
//	// 如果文件不存在则创建目录
//	if !fileExists(output) {
//		err := mkdir(output, 0755)
//		if err != nil {
//			return fmt.Errorf("preparing destination: %v", err)
//		}
//	}
//	//开始读取文件
//	file, err := os.Open(input)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//	//gzip解压
//	gr, err := gzip.NewReader(file)
//	if err != nil {
//		return err
//	}
//	defer gr.Close()
//	//tar 解压
//	//tr := tar.NewReader(gr)
//	for {
//		//err := t.untarNext(destination)
//		if err == io.EOF {
//			break
//		}
//		if err != nil {
//			//if t.ContinueOnError || IsIllegalPathError(err) {
//			//	log.Printf("[ERROR] Reading file in tar archive: %v", err)
//			//	continue
//			//}
//			return fmt.Errorf("reading file in tar archive: %v", err)
//		}
//	}
//	return nil
//}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func mkdir(dirPath string, dirMode os.FileMode) error {
	err := os.MkdirAll(dirPath, dirMode)
	if err != nil {
		return fmt.Errorf("%s: making directory: %v", dirPath, err)
	}
	return nil
}

// multipleTopLevels returns true if the paths do not
// share a common top-level folder.
func multipleTopLevels(paths []string) bool {
	if len(paths) < 2 {
		return false
	}
	var lastTop string
	for _, p := range paths {
		p = strings.TrimPrefix(strings.Replace(p, `\`, "/", -1), "/")
		for {
			next := path.Dir(p)
			if next == "." {
				break
			}
			p = next
		}
		if lastTop == "" {
			lastTop = p
		}
		if p != lastTop {
			return true
		}
	}
	return false
}

// folderNameFromFileName returns a name for a folder
// that is suitable based on the filename, which will
// be stripped of its extensions.
func folderNameFromFileName(filename string) string {
	base := filepath.Base(filename)
	firstDot := strings.Index(base, ".")
	if firstDot > -1 {
		return base[:firstDot]
	}
	return base
}
func writeNewFile(fpath string, in io.Reader, fm os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(fpath), 0755)
	if err != nil {
		return fmt.Errorf("%s: making directory for file: %v", fpath, err)
	}

	out, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("%s: creating new file: %v", fpath, err)
	}
	defer out.Close()

	err = out.Chmod(fm)
	if err != nil && runtime.GOOS != "windows" {
		return fmt.Errorf("%s: changing file mode: %v", fpath, err)
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("%s: writing file: %v", fpath, err)
	}
	return nil
}

func writeNewSymbolicLink(fpath string, target string) error {
	err := os.MkdirAll(filepath.Dir(fpath), 0755)
	if err != nil {
		return fmt.Errorf("%s: making directory for file: %v", fpath, err)
	}

	_, err = os.Lstat(fpath)
	if err == nil {
		err = os.Remove(fpath)
		if err != nil {
			return fmt.Errorf("%s: failed to unlink: %+v", fpath, err)
		}
	}

	err = os.Symlink(target, fpath)
	if err != nil {
		return fmt.Errorf("%s: making symbolic link for: %v", fpath, err)
	}
	return nil
}

func writeNewHardLink(fpath string, target string) error {
	err := os.MkdirAll(filepath.Dir(fpath), 0755)
	if err != nil {
		return fmt.Errorf("%s: making directory for file: %v", fpath, err)
	}

	_, err = os.Lstat(fpath)
	if err == nil {
		err = os.Remove(fpath)
		if err != nil {
			return fmt.Errorf("%s: failed to unlink: %+v", fpath, err)
		}
	}

	err = os.Link(target, fpath)
	if err != nil {
		return fmt.Errorf("%s: making hard link for: %v", fpath, err)
	}
	return nil
}
