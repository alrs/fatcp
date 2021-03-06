package main

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

const srcPath = "./main_test.go"

func TestCheckDirExists(t *testing.T) {
	conditions := map[string]bool{
		"./":               true,
		"/etc":             true,
		"/not/very/likely": false,
	}
	for c, b := range conditions {
		exists, err := checkDirExists(c)
		if err != nil {
			t.Fatal(err)
		}
		if exists == b {
			t.Logf("Expected: Existence of %s is %t.", c, b)
		} else {
			t.Fatalf("Failed: Existence of %s should be %t.", c, !b)
		}
	}
}

func TestCp(t *testing.T) {
	destPath := fmt.Sprint("/tmp/fatcp-test-%d", time.Now().Unix())
	err := cp(srcPath, string(destPath))
	if err != nil {
		t.Fatalf("Failed to copy file. %s", err)
	} else {
		t.Logf("Successfully copied to %s", destPath)
	}

	cleanupDestFile := func() {
		if os.Remove(destPath) == nil {
			t.Logf("Removed temporary file %s", destPath)
		} else {
			t.Fatalf("Failed to remove temporary file %s", destPath)
		}
	}
	defer cleanupDestFile()

	srcSha := sha1.New()
	destSha := sha1.New()

	srcHandle, err := os.Open(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	defer srcHandle.Close()

	destHandle, err := os.Open(destPath)
	if err != nil {
		t.Fatal(err)
	}
	defer destHandle.Close()

	if _, err := io.Copy(srcSha, srcHandle); err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(destSha, destHandle); err != nil {
		t.Fatal(err)
	}

	stringShaSum := func(shaSum hash.Hash) string {
		return fmt.Sprint("%x", shaSum.Sum(nil))
	}

	if stringShaSum(srcSha) == stringShaSum(destSha) {
		t.Log("SHA1 match between source and destination.")
	} else {
		t.Fatal("SHA1 mismatch between source and destination.")
	}
}

func TestCreateDirectory(t *testing.T) {
	testDirPath := fmt.Sprint("/tmp/fatcp-test-dir-%d", time.Now().Unix())
	defer os.Remove(testDirPath)
	err := createDirectory(testDirPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(testDirPath)
	if os.IsNotExist(err) {
		t.Fatalf("Failed to create directory %s.", testDirPath)
	}
	exists, err := checkDirExists(testDirPath)
	if err != nil {
		t.Fatal(err)
	}
	if exists == true {
		t.Logf("Expected: Directory %s created.", testDirPath)
	} else {
		t.Fatalf("Failed: Directory %s not created.", testDirPath)
	}
}

func TestSlugFilename(t *testing.T) {
	filenames := map[string]string{
		"This <is> a Terrible Filename?.mp3":   "this-is-a-terrible-filename.mp3",
		"Filename  ??  <without> an extension": "filename-without-an-extension",
	}
	for ugly, expected := range filenames {
		slugged, err := slugFileName(ugly)
		if err != nil {
			t.Fatal(err)
		}
		if slugged == expected {
			t.Logf("Successfully slugged %s", slugged)
		} else {
			t.Fatalf("%s slugged into %s", ugly, slugged)
		}
	}
}

func TestSplitPath(t *testing.T) {
	splitLength := 3
	path := "/var/log/syslog"
	split, err := splitPath(path)
	if len(split) != splitLength {
		t.Fatalf("Split path is %d elements, should be %d.", len(split), splitLength)
	} else {
		t.Logf("Split path is correctly %d elements.", splitLength)
	}
	if err != nil {
		t.Fatal(err)
	}

	expectedSlice := []string{"var", "log", "syslog"}
	if reflect.DeepEqual(split, expectedSlice) {
		t.Logf("splitPath() returned expected slice: %s", expectedSlice)
	} else {
		t.Fatalf("splitPath() returned unexpected slice: %s", expectedSlice)
	}

	badPath := "not/a/fully/qualified/path"
	_, err = splitPath(badPath)
	if err == nil {
		t.Fatalf("splitPath() failed to detect %s", badPath)
	} else {
		t.Logf("splitPath() properly detected %s", badPath)
	}

}
