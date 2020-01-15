package fs

import (
	"os"
	"strings"
	"time"
)

/*----------------FILTER 1 : AGE FILTER ------------------*/
//ageFileFilter is concrete FileFilter based on age of file
type ageFileFilter struct {
	cutoff time.Time
}

/*
NewAgeFileFilter returns ageFileFilter,
based on -
1. Time Dimension(Files older than d which can be milliseconds or seconds or minutes or hours etc.
2. multiplier of the time dimension
Example : for listing files older than 10 minutes from a directory:
param 1 = time.Minute
param 2 = 10
 */
func NewAgeFileFilter(d time.Duration, t int64) ageFileFilter{
	return ageFileFilter{
		time.Now().Add(-d* time.Duration(t)),
	}
}

func (aFf ageFileFilter) accept(f os.FileInfo) bool {
	return f.ModTime().Before(aFf.cutoff)
}

/*----------------FILTER 2 : File Name Matching FILTER ------------------
Idea is to have a simple File Name Matching filter, which simply checks if filename contains
the provided string.
This is to avoid regex matching.
*/


type fileNameFilter struct {
	name string
}

func NewNameFileFilter(name string) fileNameFilter {
	return fileNameFilter{name}
}

func (nFf fileNameFilter) accept(f os.FileInfo) bool {
	return strings.Contains(f.Name(),nFf.name)
}