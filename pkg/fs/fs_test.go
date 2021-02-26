package fs

import (
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestAgeFilter(t *testing.T){
	tmpDir := t.TempDir()
	dir1, err := ioutil.TempDir(tmpDir,"Dir1")
	dir2, err := ioutil.TempDir(tmpDir,"Dir2")
	if err != nil {
		t.Error("could not create the directory",err)
	}
	file1Dir1, _ := ioutil.TempFile(dir1,"Dir1Files_1")
	file2Dir2, _ := ioutil.TempFile(dir2,"Dir2Files_2")


	rfso := NewFileSysOp(tmpDir)

	cases := [] struct {
		name string
		durationMicroSec int64
		waitMicroSec int
		expected []string
	}{
		{"files older than given duration should be listed",
			10,
			100,
			[]string {file1Dir1.Name(),file2Dir2.Name()},
		},
		{"files not older than given duration should not be listed",
			10000,
			1,
			[]string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			aFf := NewAgeFileFilter(time.Microsecond, tc.durationMicroSec)
			time.Sleep(time.Duration(tc.waitMicroSec))
			paths, _ := rfso.ConcurrentListFiles(aFf)
			if len(tc.expected) == 0{
				if len(paths) != 0 {
					t.Fatalf("paths has %s",paths[0])
				}
			} else {
				if !(paths[0] == tc.expected[0] || paths[0] == tc.expected[1]){
					t.Fatalf("paths has %s",paths[0])
				}
				if !(paths[1] == tc.expected[0] || paths[1] == tc.expected[1]){
					t.Fatalf("paths has %s",paths[0])
				}
			}
		})
	}
}

func TestNameFilter(t *testing.T){
	tmpDir := t.TempDir()
	dir1, err := ioutil.TempDir(tmpDir,"Dir3")
	dir2, err := ioutil.TempDir(tmpDir,"Dir4")
	if err != nil {
		t.Error("could not create the directory",err)
	}
	file1Dir1, _ := ioutil.TempFile(dir1,"Dir1Files_1")
	file2Dir2, _ := ioutil.TempFile(dir2,"Dir2Files_2")

	rfso := NewFileSysOp(tmpDir)


	cases := [] struct {
		name string
		fileName string
		expected []string
	}{
		{"files with names containing given string should be listed",
			"Dir1",
			[]string {file1Dir1.Name()},
		},
		{"files with names not containing given string should not be listed",
			"Dct",
			[]string{},
		},
		{"files with names containing given string should be listed",
			"Dir2",
			[]string {file2Dir2.Name()},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			aFf := NewNameFileFilter(tc.fileName)
			paths, _ := rfso.ConcurrentListFiles(aFf)
			if len(tc.expected) == 0{
				if len(paths) != 0 {
					t.Fatalf("paths has %s",paths[0])
				}
			} else {
				if !(paths[0] == tc.expected[0]){
					t.Fatalf("paths has %s",paths[0])
				}
			}
		})
	}
}