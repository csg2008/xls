package xls

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/extrame/ole2"
)

//Open one xls file
func Open(file string, charset string) (*WorkBook, error) {
	if fi, err := os.Open(file); err == nil {
		return OpenReader(fi, charset)
	} else {
		return nil, err
	}
}

//OpenWithBuffer open one xls file with memory buffer
func OpenWithBuffer(file string, charset string) (*WorkBook, error) {
	if fi, err := ioutil.ReadFile(file); err == nil {
		return OpenReader(bytes.NewReader(fi), charset)
	} else {
		return nil, err
	}
}

//OpenWithCloser open one xls file and return the closer
func OpenWithCloser(file string, charset string) (*WorkBook, io.Closer, error) {
	if fi, err := os.Open(file); err == nil {
		wb, err := OpenReader(fi, charset)
		return wb, fi, err
	} else {
		return nil, nil, err
	}
}

//OpenReader open xls file from reader
func OpenReader(reader io.ReadSeeker, charset string) (wb *WorkBook, err error) {
	var ole *ole2.Ole
	if ole, err = ole2.Open(reader, charset); err == nil {
		var dir []*ole2.File
		if dir, err = ole.ListDir(); err == nil {
			var book *ole2.File
			var root *ole2.File
			for _, file := range dir {
				name := strings.ToLower(file.Name())
				if name == "workbook" {
					book = file
					// break
				}
				if name == "book" {
					book = file
					// break
				}
				if name == "root entry" {
					root = file
				}
			}
			if book != nil {
				wb = newWorkBookFromOle2(ole.OpenFile(book, root))
				return
			}
		}
	}
	return
}
