package main

import (
	"flag"
	"fmt"
	"github.com/extemporalgenome/slug"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// copyToFat accepts a source and destination pathname as its arguments. FAT
// filesystems have more restrictive file naming conventions than
// traditional UNIX filesystems, so copyToFat uses a slug library to
// simplify file and directory names on the destination filesystem.
func copyToFat(src string, dest string) {
	splitSrc, err := splitPath(src)
	if err != nil {
		log.Fatal(err)
	}
	rootDepth := len(splitSrc)
	walkFunc := func(path string, info os.FileInfo, err error) (localErr error) {
		if !info.IsDir() {
			splitPath, localErr := splitPath(path)
			if localErr != nil {
				return localErr
			}
			// Create a list to assemble the fully-qualified path of the destination.
			var destList []string
			// Start the list with the destination directory that was passed
			// to the function.
			destList = append(destList, dest)
			// append all of the subdirectories from the source, not including the root
			// directory or above as passed into the function, nor the filename.
			// Strip characters that will confuse a FAT filesystem.
			for _, e := range splitPath[rootDepth : len(splitPath)-1] {
				destList = append(destList, slug.Slug(e))
			}
			// Slug the filename, retain the filename extension.
			fileName, err := slugFileName(splitPath[len(splitPath)-1])
			if err != nil {
				log.Fatal(err)
			}
			// Append the slugged filename to the list representing the destination
			// directory.
			destList = append(destList, fileName)
			// Join the list into a string representing a path using the OS-specific
			// path seperator.
			destFQP := strings.Join(destList, string(os.PathSeparator))
			// This looks like an error is being thrown away, but instead it means
			// that in this case we're only interested in the filename, not the file.
			destDir, _ := filepath.Split(destFQP)
			_, err = os.Stat(destDir)
			if os.IsNotExist(err) {
				log.Printf("creating directory %s\n", destDir)
				os.MkdirAll(destDir, 0777)
			}
			err = cp(path, destFQP)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Copied %s to %s\n", path, destFQP)
		}
		return
	}
	filepath.Walk(src, walkFunc)
}

// cp copies a file from one part of the filesystem to another.
func cp(src string, dest string) error {
	srcHandle, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcHandle.Close()
	destHandle, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err := io.Copy(destHandle, srcHandle); err != nil {
		destHandle.Close()
		return err
	}
	return destHandle.Close()
}

// slugFileName takes the name of a file as its argument, and returns
// a slugged version without harming the filename extension.
func slugFileName(fn string) (string, error) {
	lastDot := strings.LastIndex(fn, ".")
	if lastDot >= 0 {
		beforeDot := slug.Slug(fn[:lastDot])
		return beforeDot + fn[lastDot:], nil
	} else {
		return slug.Slug(fn), nil
	}
}

// splitPath takes a fully-qualified path and splits the directories
// into a list.
func splitPath(path string) ([]string, error) {
	if path[0] != os.PathSeparator {
		return nil, fmt.Errorf("%s not a fully-qualified path", path)
	}
	cleanPath := filepath.Clean(string(path[1:]))
	return strings.Split(cleanPath, string(os.PathSeparator)), nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		log.Fatal("fatcp requires a source and a destination argument.")
	}

	srcDir, err := filepath.Abs(flag.Args()[0])
	if err != nil {
		log.Fatal("source directory is not a valid directory name.")
	}
	destDir, err := filepath.Abs(flag.Args()[1])
	if err != nil {
		log.Fatal("destination is not a valid directory name.")
	}

	srcStat, err := os.Stat(srcDir)
	if os.IsNotExist(err) {
		log.Fatal("source directory does not exist.")
	}
	if !srcStat.IsDir() {
		log.Fatal("source is not a directory.")
	}

	copyToFat(srcDir, destDir)
}
