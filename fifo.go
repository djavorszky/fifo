package fifo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// RecursiveCopyAll copies the source directories to the destination recursively. If sources are files, they will be
// copied to the destination folder.
//
// Returns an error if any of the sources are empty string or the destination is a file.
func RecursiveCopyAll(sources []string, destination string) error {

	return nil
}

// RecursiveCopy copies the source directory to the destination recursively. If source is a file, and destination is a
// folder, then the file will be created, otherwise renamed to the filename specified in the destination
func RecursiveCopy(src, dst string) error {
	if src == "" {
		return fmt.Errorf("src is empty")
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("src does not exist: %s", src)
		}
		return fmt.Errorf("error reading src folder %s: %v", src, err)
	}

	if dst == "" {
		return fmt.Errorf("dst is empty")
	}

	if srcInfo.IsDir() {
		return copyDirToDest(src, dst)
	}

	return copyFileToDest(src, dst)
}

func copyDirToDest(src, dst string) error {
	info, err := os.Stat(dst)
	if err != nil {
		// create dst if it does not yet exist
		if os.IsNotExist(err) {
			os.Mkdir(dst, os.ModePerm)
		} else {
			return fmt.Errorf("error reading dst folder %s: %v", dst, err)
		}
	}

	if !info.IsDir() {
		return fmt.Errorf("can't copy dir to file")
	}

	fileInfos, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("can't read source directory: %v", err)
	}

	for _, fi := range fileInfos {
		if fi.IsDir() {
			err = copyDirToDest(fi.Name(), filepath.Join(dst, fi.Name()))
		} else {
			err = copyFileToDest(fi.Name(), filepath.Join(dst, fi.Name()))
		}

		if err != nil {
			return fmt.Errorf("failed copying file %s: %v", fi.Name(), err)
		}
	}

	return nil
}

func copyFileToDest(src, dst string) error {

	return nil
}

// fileName checks whether the path points to a filename or not. Hidden files/folders are determined to be folders.
func fileName(path string) bool {
	_, f := filepath.Split(path)
	if strings.Contains(f, ".") && strings.Index(f, ".") != 0 {
		return true
	}

	return false
}
