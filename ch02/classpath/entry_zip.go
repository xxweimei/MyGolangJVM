package classpath

import (
	"path/filepath"
	"archive/zip"
	"errors"
	"io/ioutil"
)

type ZipEntry struct {
	absPath string
	zipRc   *zip.ReadCloser
}

func newZipEntry(path string) *ZipEntry {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return &ZipEntry{absPath, nil}
}

func (self *ZipEntry) readClass(className string) ([]byte, Entry, error) {
	if self.zipRc == nil {
		err := self.openJar()
		defer self.closeJar()
		if err != nil {
			return nil, nil, err
		}
	}

	classFile := self.findClass(className)
	if classFile == nil {
		return nil, nil, errors.New("class not found: " + className)
	}

	rc, err := classFile.Open()
	defer rc.Close()
	if err != nil {
		return nil, nil, err
	}
	// read class data
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, nil, err
	}
	return data, self, nil
}

func (self *ZipEntry) String() string {
	return self.absPath
}

func (self *ZipEntry) openJar() error {
	r, err := zip.OpenReader(self.absPath)
	if err == nil {
		self.zipRc = r
	}

	return err
}

func (self *ZipEntry) closeJar() {
	if self.zipRc != nil {
		self.zipRc.Close()
	}
}

func (self *ZipEntry) findClass(className string) *zip.File {
	for _, f := range self.zipRc.File {
		if f.Name == className {
			return f
		}
	}
	return nil
}
