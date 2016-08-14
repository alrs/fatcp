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

func TestCp(t *testing.T) {
	destPath := fmt.Sprintf("/tmp/fatcp-test-%d", time.Now().Unix())
	err := cp(srcPath, string(destPath))
	if err != nil {
		t.Fatalf("Failed to copy file. %s", err)
	} else {
		t.Logf("Successfully copied to %s", destPath)
	}

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
