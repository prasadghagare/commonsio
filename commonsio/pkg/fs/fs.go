package fs

import (
	"os"
)

//FileSysOp enhances the File object from stdlib
type FileSysOp interface {
	//open(path string) (FileInfo, error)
	ListFiles(FileFilter) ([]string, error)
	IsDir() bool
}

//FileFilter : anything that implement accept
type FileFilter interface {
	accept(os.FileInfo) bool
}

