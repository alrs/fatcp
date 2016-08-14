package main

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"os"
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
