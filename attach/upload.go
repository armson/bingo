package attach

//import (
//	"path/filepath"
//	"os"
//	"io"
//	"errors"
//	"github.com/armson/bingo/config"
//	"fmt"
//	"strings"
//)

//func (img *Image)SetPath(path string) *Image {
//	img.path = path
//	return img
//}
//
//func (img *Image)SetName(fileName string) *Image {
//	img.fileName = fileName
//	return img
//}
//
//func (img *Image)Save() (bool , error) {
//	if _, err := img.checkSize(); err != nil {
//		return false, err
//	}
//
//	if _, err := img.checkMime(); err != nil {
//		return false, err
//	}
//
//	path, err := filepath.Abs(img.path)
//	if err != nil { return false, err }
//
//	if img.fileName == "" { img.fileName = img.Name() }
//	fileName := filepath.Join(path, img.fileName)
//
//	os.MkdirAll(filepath.Dir(fileName), 0777)
//
//	fw , err := os.Create(fileName)
//	if err != nil {  return false, err }
//	defer fw.Close()
//
//	_, err = io.Copy(fw, img.File)
//	if err != nil { return  false , err }
//
//	return true , nil
//}
//
//func (img *Image)checkSize() (bool , error) {
//	size := img.Size()
//	if size == 0 {
//		return false , errors.New("File size must be greater than zero")
//	}
//
//	uploadMaxSize  := config.Int("uploadmaxsize")
//	if size > uploadMaxSize {
//		return false , fmt.Errorf("File size must be less than %dM" , uploadMaxSize/1024/1024)
//	}
//	return  true ,nil
//}
//
//func (img *Image)checkMime() (bool , error) {
//	uploadPermitExt  := config.Slice("uploadPermitExt")
//	if len(uploadPermitExt) < 1 {
//		return true, nil
//	}
//
//	ext  := img.Ext()
//	mime := img.MimeType()
//	for _, v := range uploadPermitExt {
//		if m , ok := Mimes[v]; ok && m == mime && v == ext {
//			return true ,nil
//		}
//	}
//	return false, fmt.Errorf("File ext must be [%s]" , strings.Join(uploadPermitExt, " "))
//}




