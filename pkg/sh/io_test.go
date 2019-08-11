package sh_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/craiggwilson/goke/pkg/sh"

	"testing"
)

func TestCopy_File(t *testing.T) {
	ctx := makeTestContext()

	tempDir, err := ioutil.TempDir("", "copyFile")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	fileAPath := filepath.Join(tempDir, "fileA")
	fileA, err := os.Create(fileAPath)
	if err != nil {
		t.Fatalf("failed creating fileA: %v", err)
	}

	fmt.Fprint(fileA, "not empty")

	_ = fileA.Close()

	fileBPath := filepath.Join(tempDir, "fileB")
	err = sh.Copy(ctx, fileAPath, fileBPath)
	if err != nil {
		t.Fatalf("failed copying file: %v", err)
	}

	fileBContents, err := ioutil.ReadFile(fileBPath)
	if err != nil {
		t.Fatalf("failed getting fileB contents: %v", err)
	}

	if string(fileBContents) != "not empty" {
		t.Fatalf("expected \"not empty\", but got %q", fileBContents)
	}
}

func TestCreateDirectory(t *testing.T) {
	ctx := makeTestContext()

	tempDir, err := ioutil.TempDir("", "createDirectory")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	testDir := filepath.Join(tempDir, "level1")
	if _, err = os.Stat(testDir); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed ensuring directory does not exist: %v", err)
	}

	err = sh.CreateDirectory(ctx, testDir)
	if err != nil {
		t.Fatalf("failed creating directory: %v", err)
	}

	if _, err = os.Stat(testDir); err != nil {
		t.Fatalf("failed ensuring directory exists: %v", err)
	}
}

func TestCreateDirectoryR(t *testing.T) {
	ctx := makeTestContext()

	tempDir, err := ioutil.TempDir("", "createDirectoryR")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	testDir := filepath.Join(tempDir, "level1", "level2")
	if _, err = os.Stat(testDir); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed ensuring directory does not exist: %v", err)
	}

	err = sh.CreateDirectoryR(ctx, testDir)
	if err != nil {
		t.Fatalf("failed creating directory: %v", err)
	}

	if _, err = os.Stat(testDir); err != nil {
		t.Fatalf("failed ensuring directory exists: %v", err)
	}
}

func TestCreateFile(t *testing.T) {
	ctx := makeTestContext()

	tempDir, err := ioutil.TempDir("", "createFile")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	testFile := filepath.Join(tempDir, "file")
	if _, err = os.Stat(testFile); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed ensuring file does not exist: %v", err)
	}

	_, err = sh.CreateFile(ctx, testFile)
	if err != nil {
		t.Fatalf("failed creating file: %v", err)
	}

	if _, err = os.Stat(testFile); err != nil {
		t.Fatalf("failed ensuring file exists: %v", err)
	}
}

func TestCreateFileR(t *testing.T) {
	ctx := makeTestContext()

	tempDir, err := ioutil.TempDir("", "createFile")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	testFile := filepath.Join(tempDir, "level1", "file")
	if _, err = os.Stat(testFile); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed ensuring file does not exist: %v", err)
	}

	_, err = sh.CreateFileR(ctx, testFile)
	if err != nil {
		t.Fatalf("failed creating file: %v", err)
	}

	if _, err = os.Stat(testFile); err != nil {
		t.Fatalf("failed ensuring file exists: %v", err)
	}
}

func TestDirectoryExists(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "directoryExists")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	exists, err := sh.DirectoryExists(tempDir)
	if err != nil {
		t.Fatalf("failed testing for directory existance: %v", err)
	}

	if !exists {
		t.Fatal("expected directory to exist")
	}

	exists, err = sh.DirectoryExists(filepath.Join(tempDir, "noexisty"))
	if err != nil {
		t.Fatalf("failed testing for directory existance: %v", err)
	}

	if exists {
		t.Fatal("expected directory not to exist")
	}
}

func TestFileExists(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "fileExists")
	if err != nil {
		t.Fatalf("failed making temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	exists, err := sh.FileExists(tempFile.Name())
	if err != nil {
		t.Fatalf("failed testing for directory existance: %v", err)
	}

	if !exists {
		t.Fatal("expected file to exist")
	}

	exists, err = sh.FileExists(filepath.Join(tempFile.Name(), "noexisty"))
	if err != nil {
		t.Fatalf("failed testing for directory existance: %v", err)
	}

	if exists {
		t.Fatal("expected file not to exist")
	}
}

func TestIsDirectoryEmpty(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "isDirectoryEmpty")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	empty, err := sh.IsDirectoryEmpty(tempDir)
	if err != nil {
		t.Fatalf("failed testing if the directory was empty: %v", err)
	}

	if !empty {
		t.Fatal("expected directory to be empty")
	}

	testFile := filepath.Join(tempDir, "file")
	_, err = os.Create(testFile)
	if err != nil {
		t.Fatalf("failed creating file: %v", err)
	}

	empty, err = sh.IsDirectoryEmpty(tempDir)
	if err != nil {
		t.Fatalf("failed testing if the directory was empty: %v", err)
	}

	if empty {
		t.Fatal("expected directory to not be empty")
	}
}

func TestIsFileEmpty(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "isFileEmpty")
	if err != nil {
		t.Fatalf("failed making temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	empty, err := sh.IsFileEmpty(tempFile.Name())
	if err != nil {
		t.Fatalf("failed testing if the file was empty: %v", err)
	}

	if !empty {
		t.Fatal("expected file to be empty")
	}

	fmt.Fprintln(tempFile, "not empty")

	empty, err = sh.IsFileEmpty(tempFile.Name())
	if err != nil {
		t.Fatalf("failed testing if the file was empty: %v", err)
	}

	if empty {
		t.Fatal("expected file to not be empty")
	}
}

func TestRemove_Directory(t *testing.T) {
	ctx := makeTestContext()

	tempDir, err := ioutil.TempDir("", "removeDirectory")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempDir)

	if _, err = os.Stat(tempDir); err != nil {
		t.Fatalf("failed ensuring directory exists: %v", err)
	}

	err = sh.Remove(ctx, tempDir)
	if err != nil {
		t.Fatalf("failed removing directory: %v", err)
	}

	if _, err = os.Stat(tempDir); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed ensuring directory does not exist: %v", err)
	}
}

func TestRemove_File(t *testing.T) {
	ctx := makeTestContext()

	tempFile, err := ioutil.TempFile("", "removeFile")
	if err != nil {
		t.Fatalf("failed making temp directory: %v", err)
	}
	defer os.Remove(tempFile.Name())
	_ = tempFile.Close()

	if _, err = os.Stat(tempFile.Name()); err != nil {
		t.Fatalf("failed ensuring file exists: %v", err)
	}

	err = sh.Remove(ctx, tempFile.Name())
	if err != nil {
		t.Fatalf("failed removing file: %v", err)
	}

	if _, err = os.Stat(tempFile.Name()); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed ensuring file does not exist: %v", err)
	}
}
