package sx_test

import (
	"os"
	"testing"

	"github.com/ZeroBsd/sx"
)

func TestFileNew(t *testing.T) {
	var tempDir = sx.DirTemp()
	if !tempDir.Exists() {
		t.FailNow()
	}
	var emptyDir = sx.NewDirFromString("")
	if emptyDir != sx.DirRoot() {
		t.FailNow()
	}

	var dir = sx.DirUserHome().Value()
	if !dir.Exists() {
		t.FailNow()
	}
	var newTempDir = ".sx_test"
	var result = dir.CreateDir(newTempDir)
	if !result.Ok() {
		t.FailNow()
	}
	var dir2 = result.Value()
	if !dir2.Exists() || !dir.IsDirectory(newTempDir) || dir.IsFile(newTempDir) {
		t.FailNow()
	}
	result = dir.CreateDir(newTempDir)
	if !result.Ok() {
		t.FailNow()
	}
	var dir3 = result.Value()
	if dir3.String() != dir2.String() {
		t.FailNow()
	}
	if !dir.HasParent() || dir.Parent().Value() != sx.DirUserHome().Value() {
		t.FailNow()
	}
	result = dir.Cd("_folder_that_does_not_exist")
	if result.Ok() {
		t.FailNow()
	}
	dir = dir.Cd(newTempDir).Value()
	if dir.String() != dir3.String() || !dir.Exists() {
		t.FailNow()
	}
	var testFileName = "sx_testfile"
	var testFileContent = "test test\r\ntest\n"
	var testFileContentClean = "test test\ntest\n"
	dir.WriteAllText(testFileName, testFileContent)
	var bytes = dir.ReadAllBytes(testFileName).Value()
	if string(bytes) != testFileContent {
		t.FailNow()
	}
	var text = dir.ReadAllText(testFileName).Value()
	if text != testFileContentClean {
		t.FailNow()
	}

	var failedDir = sx.NewDirFromString("/::")
	if e := failedDir.CreateDir("this fails, because dir cannot exist"); e.Ok() {
		t.FailNow()
	}

	var found = false
	for it := dir.NewIterator(); it.Ok(); it.Next() {
		if it.Value() != "sx_testfile" {
			t.FailNow()
		} else {
			found = true
			continue
		}
		t.FailNow()
	}
	if !found {
		t.FailNow()
	}

}

func TestFailedFiles(t *testing.T) {
	var dir = sx.DirUserHome().Value()
	if dir.IsSymlink("_link_that_does_not_exist") {
		t.FailNow()
	}
	if e := dir.ReadAllBytes("_file_that_does_not_exist"); e.Ok() {
		t.FailNow()
	}
	if e := dir.ReadAllText("_file_that_does_not_exist"); e.Ok() {
		t.FailNow()
	}
	if dir := sx.NewDir("C:"); !dir.Exists() || dir.String() != sx.DirRoot().String() {
		t.FailNow()
	}
}

func TestTemp(t *testing.T) {
	var mp = sx.DirTemp()
	var dir = sx.MkDirAndPath(mp.String() + "/.sx_test/someFolder")
	if !dir.Ok() {
		t.FailNow()
	}
	var failed = sx.MkDirAndPath("NOTEXISTING:/.sx_test")
	if failed.Ok() {
		t.FailNow()
	}
}

func TestDirRoot(t *testing.T) {
	var dirRoot = sx.DirRoot()
	if dir := sx.NewDir(); dir != dirRoot {
		t.FailNow()
	}
	if dir := sx.NewDir(""); dir != dirRoot {
		t.FailNow()
	}
}

func TestFailedDirHome(t *testing.T) {
	os.Unsetenv("USERPROFILE")
	if dir := sx.DirUserHome(); dir.Ok() {
		t.FailNow()
	}
}

func TestFileBaseAndExtension(t *testing.T) {
	if ext := sx.FileExtension("someFile.txt"); ext != ".txt" {
		t.FailNow()
	}
	if ext := sx.FileExtension("somePath/someFile.txt"); ext != ".txt" {
		t.FailNow()
	}
	if ext := sx.FileExtension("someFile.bat.txt"); ext != ".txt" {
		t.FailNow()
	}
	if base := sx.FileBaseName("somePath/someFile.txt"); base != "someFile.txt" {
		t.FailNow()
	}
	if base := sx.FileBaseName("somePath\\someFile.tar.gz"); base != "someFile.tar.gz" {
		t.FailNow()
	}
	if noext := sx.FileWithoutExtension("somePath\\someFile.tar.gz"); noext != "someFile.tar" {
		t.FailNow()
	}
	if noext := sx.FileWithoutExtension(""); noext != "" {
		t.FailNow()
	}
}
