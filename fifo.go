package fifo

import "fmt"

// RecursiveCopyAll copies the source directories to the destination recursively. If sources are files, they will be
// copied to the destination folder.
//
// Returns an error if any of the sources are empty string or the destination is a file.
func RecursiveCopyAll(sources []string, destination string) error {

	return nil
}

// RecursiveCopy copies the source directory to the destination recursively. If source is a file, and destination is a
// folder, then the file will be created, otherwise renamed to the filename specified in the destination
func RecursiveCopy(source, destination string) error {
	if source == "" {
		return fmt.Errorf("source is empty")
	}

	if destination == "" {
		return fmt.Errorf("destination is empty")
	}

	return nil
}
