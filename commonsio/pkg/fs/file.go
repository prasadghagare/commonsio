package fs

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

//NewFileSysOp gives RealFileSysOp struct instance
//This instance implements FileSysOp interface,
//and hence provides client various filesystem ops beyond stdlib
//Use the methods implemented on this with filters provided.
func NewFileSysOp(path string) FileSysOp {
	f, err := os.Open(path)
	if err != nil {
		panic("error opening file")
	}
	return realFile {
		File: f,
	}
}

type realFile struct {
	*os.File
}

func (f realFile) IsDir() bool {
	fstat,_ := f.File.Stat()
	return fstat.IsDir()
}

/*
ListFiles returns a slice of string containing files' name(not directories) based on the Filter used.
 */
func (f realFile) ListFiles(fF FileFilter) (matchedFiles []string, err error) {
	if f.IsDir() {
		err = filepath.Walk(f.Name(),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return errors.Wrapf(err, "Could not list for path %q: ", path)
				}
			if (fF.accept(info)) && ! info.IsDir() {
				matchedFiles = append(matchedFiles, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		return nil , errors.Errorf("%q is a file, ListFiles needs a directory",f.Name())
	}

	return
}