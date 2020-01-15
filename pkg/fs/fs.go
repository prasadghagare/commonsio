package fs

import (
	"os"
)

//FileSysOp enhances the File object from stdlib
type FileSysOp interface {
	ListFiles(FileFilter) ([]string, error)
	ConcurrentListFiles(FileFilter)([]string, error)
	IsDir() bool
}

//FileFilter : anything that implement accept
type FileFilter interface {
	accept(os.FileInfo) bool
}

