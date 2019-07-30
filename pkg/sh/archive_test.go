package sh_test

import (
	"os"
	"testing"

	"github.com/craiggwilson/goke/pkg/sh"
)

func TestArchiveZIP(t *testing.T) {
	ctx := makeTestContext()

	err := sh.Archive(ctx, "./testdata/level1", "./testdata/temp.zip")
	if err != nil {
		t.Fatalf("failed archiving: %v", err)
	}
	defer os.Remove("./testdata/temp.zip")

	err = sh.Unarchive(ctx, "./testdata/temp.zip", "./testdata/level1-zip")
	if err != nil {
		t.Fatalf("failed unarchiving: %v", err)
	}
	defer os.RemoveAll("./testdata/level1-zip")
}

func TestArchiveTGZ(t *testing.T) {
	ctx := makeTestContext()

	err := sh.Archive(ctx, "./testdata/level1", "./testdata/temp.tgz")
	if err != nil {
		t.Fatalf("failed archiving: %v", err)
	}
	defer os.Remove("./testdata/temp.tgz")

	err = sh.Unarchive(ctx, "./testdata/temp.tgz", "./testdata/level1-tgz")
	if err != nil {
		t.Fatalf("failed unarchiving: %v", err)
	}
	defer os.RemoveAll("./testdata/level1-tgz")
}
