package attach

import (
	"mime/multipart"
	"path"
	"os"
	"io"
	"path/filepath"
	"strings"
)

type Attachment struct {
	file multipart.File
	header *multipart.FileHeader

	Path string
	Name string
	Ext string
	Mime string
	Size int
}

func New(file multipart.File, header *multipart.FileHeader) *Attachment {
	a := &Attachment{file:file, header:header}
	a.Name = a.name()
	a.Ext  = a.ext()
	a.Mime = a.mimeType()
	a.Size = a.size()
	return a
}

func (a *Attachment)ext() string {
	return strings.ToLower(path.Ext(a.header.Filename))
}
func (a *Attachment)mimeType() string {
	if a.header.Header.Get("Content-Type") == "image/jpg" {
		a.header.Header.Set("Content-Type", "image/jpeg")
	}
	return strings.ToLower(a.header.Header.Get("Content-Type"))
}

func (a *Attachment)name() string {
	if a.header.Filename != "" {
		_, fileName := filepath.Split(a.header.Filename)
		return strings.ToLower(fileName)
	}
	return ""
}

func (a *Attachment)Rc() io.ReadCloser {
	return  a.file.(io.ReadCloser)
}

type sizer interface {
	Size() int64
}

type stater interface {
	Stat() (os.FileInfo, error)
}

func (a *Attachment)size() int {
	if statInterface, ok := a.file.(stater); ok {
		fileInfo, _ := statInterface.Stat()
		return int(fileInfo.Size())
	}
	if sizeInterface, ok := a.file.(sizer); ok {
		return int(sizeInterface.Size())
	}
	return 0
}

func (a *Attachment)Header() *multipart.FileHeader {
	return a.header
}


