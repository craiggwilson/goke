package sh

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/craiggwilson/goke/task"
)

// Archive will create an archive from the src file or directory and use the destination's
// extension to determine which format to use.
func Archive(ctx *task.Context, src, dest string) error {
	if strings.HasSuffix(dest, ".zip") {
		return ArchiveZip(ctx, src, dest)
	}

	return errors.New("unable to determine archive format")
}

// ArchiveZip will zip the src into a a zipped file at the dest.
func ArchiveZip(ctx *task.Context, src, dest string) error {
	ctx.Logf("zip: %s -> %s\n", src, dest)

	src, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	dest, err = filepath.Abs(dest)
	if err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	var baseDir string
	if srcFileInfo.IsDir() {
		baseDir = filepath.Base(src)
	}

	filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == src || path == dest {
			return nil
		}

		header, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, src))
		}

		if fi.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		destFile, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})

	return nil
}

// Unarchive decompresses the archive according to the source's extension.
func Unarchive(ctx *task.Context, src, dest string) error {
	if strings.HasSuffix(src, ".zip") {
		return UnarchiveZip(ctx, src, dest)
	}

	return errors.New("unable to determine archive format")
}

// UnarchiveZip decompresses the src zip file into the destination.
func UnarchiveZip(ctx *task.Context, src, dest string) error {
	ctx.Logf("unzip: %s -> %s\n", src, dest)

	src, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	dest, err = filepath.Abs(dest)
	if err != nil {
		return err
	}

	zr, err := zip.OpenReader(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, file := range zr.File {
		path := filepath.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		srcFile, err := file.Open()
		if err != nil {
			return err
		}

		defer srcFile.Close()

		destFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer destFile.Close()

		if _, err = io.Copy(destFile, srcFile); err != nil {
			return err
		}
	}

	return nil
}
