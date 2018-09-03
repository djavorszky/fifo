package fifo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	testDirs = []string{
		filepath.Join("one", "a", "temp"),
		filepath.Join("one", "b"),
		filepath.Join("two", "c", "notPorn"),
	}
	testFileLocations = []string{
		filepath.Join("one"),
		filepath.Join("one", "a"),
		filepath.Join("one", "a", "temp"),
		filepath.Join("one", "b"),
		filepath.Join("two"),
		filepath.Join("two"),
		filepath.Join("two", "c", "notPorn"),
	}

	testSrc, testDst string
)

func setup() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed getting current dir: %v", err)
	}

	testSrc = filepath.Join(dir, "testSrc")
	testDst = filepath.Join(dir, "testDst")

	err = os.Mkdir(testSrc, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed creating test source directory %s: %v", testSrc, err)
	}

	err = os.Mkdir(testDst, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed creating test destination directory %s: %v", testDst, err)
	}

	for _, d := range testDirs {
		err := os.MkdirAll(filepath.Join(testSrc, d), os.ModePerm)
		if err != nil {
			// cleanup
			os.RemoveAll(testSrc)
			return fmt.Errorf("failed creating directories: %v", err)
		}
	}

	for _, d := range testFileLocations {
		f, err := ioutil.TempFile(filepath.Join(testSrc, d), "tempfile")
		if err != nil {
			// cleanup
			os.RemoveAll(testSrc)
			return fmt.Errorf("failed creating tempfiles: %v", err)
		}
		defer f.Close()

		// Write the file's path into the file
		_, err = f.WriteString(f.Name())
		if err != nil {
			return fmt.Errorf("failed writing text into tempfile: %v", err)
		}
	}

	return nil
}

func teardown() error {
	err := os.RemoveAll(testSrc)
	if err != nil {
		return fmt.Errorf("failed removing test source directory: %v", err)
	}

	err = os.RemoveAll(testDst)
	if err != nil {
		return fmt.Errorf("failed removing test destination directory: %v", err)
	}

	return nil
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		fmt.Printf("Failed setup: %v", err)
		os.Exit(-1)
	}

	res := m.Run()

	err = teardown()
	if err != nil {
		fmt.Printf("Failed teardown: %v", err)
		os.Exit(-1)
	}

	os.Exit(res)
}

func TestRecursiveCopy(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty_source", args{"", "/some/destination"}, true},
		{"empty_destination", args{"/some/source", ""}, true},
		{"dir_with_files", args{"/some/source", ""}, true},
		{"empty_destination", args{"/some/source", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RecursiveCopy(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("RecursiveCopy() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}
		})
	}
}
