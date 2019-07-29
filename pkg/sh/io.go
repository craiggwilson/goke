package sh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/craiggwilson/goke/task"
)

// CopyFile copies a file.
func CopyFile(ctx *task.Context, fromPath, toPath string) error {
	ctx.Logf("copy file: %s -> %s\n", fromPath, toPath)
	from, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed opening %s: %v", fromPath, err)
	}
	defer from.Close()

	fi, err := from.Stat()
	if err != nil {
		return fmt.Errorf("failed statting %s: %v", fromPath, err)
	}

	to, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed creating/opening %s: %v", toPath, err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return fmt.Errorf("failed copying %s to %s: %v", fromPath, toPath, err)
	}

	return nil
}

// CreateDirectory creates a directory.
func CreateDirectory(ctx *task.Context, path string) error {
	ctx.Logf("create dir: %s\n", path)
	err := os.Mkdir(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed making directory %s: %v", path, err)
	}

	return nil
}

// CreateDirectoryR creates a directory recursively.
func CreateDirectoryR(ctx *task.Context, path string) error {
	ctx.Logf("create dir recursive: %s\n", path)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed making directory %s: %v", path, err)
	}

	return nil
}

// CreateFile creates a file.
func CreateFile(ctx *task.Context, path string) (*os.File, error) {
	ctx.Logf("create file: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed creating file %s: %v", path, err)
	}

	return f, nil
}

// CreateFileR creates a file ensuring all the directories are created recursively.
func CreateFileR(ctx *task.Context, path string) (*os.File, error) {
	ctx.Logf("create file recursive: %s\n", path)
	dir := filepath.Dir(path)

	err := CreateDirectoryR(ctx, dir)
	if err != nil {
		return nil, err // already has a good error message
	}

	return CreateFile(ctx, path)
}

// DirectoryExists indicates if the directory exists.
func DirectoryExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("failed statting path %s: %v", path, err)
	}

	return fi.IsDir(), nil
}

// FileExists indicates if the file exists.
func FileExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, fmt.Errorf("failed statting path %s: %v", path, err)
	}

	return !fi.IsDir(), nil
}

// IsDirectoryEmpty indicates if the directory is empty.
func IsDirectoryEmpty(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("directory %s does not exist", path)
		}

		return false, fmt.Errorf("failed statting path %s: %v", path, err)
	}

	if !fi.IsDir() {
		return false, fmt.Errorf("%s is not a directory", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("failed opening %s: %v", path, err)
	}
	defer f.Close()
	entries, _ := f.Readdir(-1)
	return len(entries) == 0, nil
}

// IsFileEmpty indicates if the file is empty.
func IsFileEmpty(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, fmt.Errorf("file %s does not exist", path)
		}

		return false, fmt.Errorf("failed statting path %s: %v", path, err)
	}

	if fi.IsDir() {
		return false, fmt.Errorf("%s is not a file", path)
	}

	return fi.Size() == 0, nil
}

// RemoveDirectory removes the directory.
func RemoveDirectory(ctx *task.Context, path string) error {
	ctx.Logf("remove directory: %s\n", path)
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	err = os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed removing %s: %v", path, err)
	}

	return nil
}

// RemoveFile removes the file.
func RemoveFile(ctx *task.Context, path string) error {
	ctx.Logf("remove file: %s\n", path)
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("%s is not a file", path)
	}

	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed removing %s: %v", path, err)
	}

	return nil
}
