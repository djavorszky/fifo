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
		"one", "two", "two",
		filepath.Join("one", "a"),
		filepath.Join("one", "a", "temp"),
		filepath.Join("one", "b"),
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
		{"dir_with_files", args{filepath.Join(testSrc, "two", "c", "notPorn"), filepath.Join(testDst, "first")}, false},
		{"empty_destination", args{"/some/source", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RecursiveCopy(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("RecursiveCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if err := srcEqualsDst(tt.args.src, tt.args.dst); err != nil {
				t.Errorf("RecursiveCopy() fail: %v", err)
			}
		})
	}

}

func srcEqualsDst(src, dst string) error {
	err := filepath.Walk(dst, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("copy %v ->  %v: %v", src, dst, err)
	}

	return nil
}

func Test_fileName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"file_linux", args{"/test/folder/file.go"}, true},
		{"file_win", args{"C:\\test\\folder\\file.go"}, true},
		{"folder_linux", args{"/test/folder/another"}, false},
		{"folder_win", args{"C:\\test\\folder\\another"}, false},
		{"folder_linux_trailing_slash", args{"/test/folder/another/"}, false},
		{"folder_win_trailing_slash", args{"C:\\test\\folder\\another\\"}, false},
		{"folder_linux_trailing_slash_with_period", args{"/test/folder/another.go/"}, false},
		{"folder_win_trailing_slash_with_period", args{"C:\\test\\folder\\another.go\\"}, false},
		{"file_linux_hidden", args{"/test/folder/.hello.go"}, false},
		{"folder_linux_hidden", args{"/test/folder/.git"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileName(tt.args.path); got != tt.want {
				t.Errorf("fileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
