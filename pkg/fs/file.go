package fs

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
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

//feedDirs is generator to set up the pipeline by putting list of dirs in root on channel
func feedDirs(outDirs []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, dir := range outDirs{
			out <- dir
		}
		close(out)
	}()
	return out
}

//readDirsInRoot - to read the all the directories in root dir passed
func readDirsInRoot(dirname string) ([]string, error){
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	var newNames []string
	for _,name := range names {
		newNames = append(newNames, filepath.Join(dirname,name))
	}
	return newNames, nil
}

func singleDirListFiles(in <-chan string, fF FileFilter) <-chan string {
	channelMatchedFileOut := make(chan string)
	err := errors.New("some error")
	go func() {
		for  n := range in {
			err = filepath.Walk(n,
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return errors.Wrapf(err, "Could not list for path %q: ", path)
					}
					if (fF.accept(info)) && ! info.IsDir() {
						channelMatchedFileOut <- path
					}
					return nil
				})
			//TODO - improve this error handling
			//This is a library which should not be logging like this rather should return the error
			//For clients to handle
			if err != nil {
				fmt.Println(err)
			}
		}
		close(channelMatchedFileOut)
	}()
	return channelMatchedFileOut
	}

func merge(cs []<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (f realFile) ConcurrentListFiles(fF FileFilter) (matchedFiles[] string, err error) {
	dirsInRoot , errListingRootDir := readDirsInRoot(f.Name())
	if errListingRootDir != nil {
		return nil, errListingRootDir
	}
	cReadDirs := feedDirs(dirsInRoot)
	numOfCPU := runtime.NumCPU()
	chans := make([]<-chan string, numOfCPU)
	for i:=0; i < numOfCPU ;i++ {
		chans[i] = singleDirListFiles(cReadDirs,fF)
	}
	for n := range merge(chans) {
		matchedFiles = append(matchedFiles,n )
	}
	return
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
