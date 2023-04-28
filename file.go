// SPDX-License-Identifier: 0BSD
package sx

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const FilePathSeparator = string(os.PathSeparator)
const RunningOnWindows = runtime.GOOS == "windows"

var _ Iterable[int, string] = NewDir()

func DirRoot() Dir {
	var root = "/"
	if RunningOnWindows {
		root = "C:\\"
	}
	return NewDirFromString(root)
}

func DirTemp() Dir {
	return NewDirFromString(os.TempDir()).normalize()
}

func DirUserHome() Result[Dir] {
	var userHomeDir, err = os.UserHomeDir()
	if err != nil {
		return NewResultFromError[Dir](err)
	}
	return NewResultFrom(NewDirFromString(userHomeDir))
}

// Returns a Dir from a string, creates a new directory if needed (including all parents)
func MkDirAndPath(fullPathString string) Result[Dir] {
	var err = os.MkdirAll(fullPathString, 700)
	if err != nil {
		return NewResultFromError[Dir](err)
	}
	return NewResultFrom(NewDirFromString(fullPathString))
}

type Dir string

// Create a Dir from a path/string
// Directories always end with the FilePathSeparator
func NewDirFromString(path string) Dir {
	path, _ = strings.CutPrefix(path, "file://")
	if path == "" {
		return DirRoot()
	}
	if RunningOnWindows {
		path = strings.ReplaceAll(path, `/`, `\`)
	}
	return Dir(path).normalize()
}

func (dir Dir) normalize() Dir {
	var path = string(dir)
	if !strings.HasSuffix(path, FilePathSeparator) {
		path = StrCat(path, FilePathSeparator)
	}
	return Dir(path)
}

// Create a Dir from its parts
// It is assumed that there are no FilePathSeparators within the parts
// Directories always end with the FilePathSeparator
func NewDir(pathParts ...string) Dir {
	if len(pathParts) == 0 {
		return DirRoot()
	}
	var path = StrJoin(FilePathSeparator, pathParts...)
	return NewDirFromString(path)
}

func (dir Dir) String() string {
	return string(dir)
}

func (dir Dir) Exists() bool {
	fileInfo, err := os.Stat(dir.String())
	return err == nil && fileInfo.Mode().IsDir()
}

func (dir Dir) IsFile(fileName string) bool {
	var fullName = StrCat(dir.String(), fileName)
	fileInfo, err := os.Stat(fullName)
	return err == nil && fileInfo.Mode().IsRegular()
}

func (dir Dir) IsDirectory(fileName string) bool {
	var fullName = StrCat(dir.String(), fileName)
	fileInfo, err := os.Stat(fullName)
	return err == nil && fileInfo.Mode().IsDir()
}

func (dir Dir) IsSymlink(fileName string) bool {
	var fullName = StrCat(dir.String(), fileName)
	fileInfo, err := os.Stat(fullName)
	return err == nil && (fileInfo.Mode()&os.ModeSymlink) != 0
}

func (dir Dir) CreateDir(folderName string) Result[Dir] {
	var newDirName = StrCat(dir.String(), folderName)
	var newDir = NewDirFromString(newDirName).normalize()
	var err = os.MkdirAll(newDir.String(), 700)
	if err != nil {
		return NewResultFromError[Dir](err)
	}
	return NewResultFrom(newDir)
}

func (dir Dir) ReadAllBytes(fromFileName string) Result[[]byte] {
	var content, err = os.ReadFile(StrCat(dir.String(), fromFileName))
	if err != nil {
		return NewResultFromError[[]byte](err)
	}
	return NewResultFrom(content)
}

// Reads all the text from a file, newlines are always converted to '\n'
func (dir Dir) ReadAllText(fromFileName string) Result[string] {
	var bytes = dir.ReadAllBytes(fromFileName)
	if !bytes.Ok() {
		return NewResultError[string](bytes.Error())
	}
	var text = string(bytes.Value())
	if RunningOnWindows {
		text = strings.ReplaceAll(text, "\r\n", "\n")
	}
	return NewResultFrom(text)
}

// Write all the text to a file, no modifications
func (dir Dir) WriteAllText(toFileName string, text string) error {
	var fullFileName = StrCat(dir.String(), toFileName)
	return os.WriteFile(fullFileName, []byte(text), fs.ModePerm)
}

// Descends into a folder. Folder must exist.
func (dir Dir) Cd(folderName string) Result[Dir] {
	if !dir.IsDirectory(folderName) {
		return NewResultError[Dir](ReflectFunctionName(), ": cannot change directory, because directory '", folderName, "' does not exist in '", dir.String(), "'")
	}
	var newDirName = StrCat(dir.String(), folderName)
	var newDir = NewDirFromString(newDirName).normalize()
	return NewResultFrom(newDir)
}

func (dir Dir) HasParent() bool {
	var parent = dir.Parent()
	return parent.Ok()
}
func (dir Dir) Parent() Result[Dir] {
	var parentDirString = filepath.Dir(dir.String())
	var parent = NewDirFromString(parentDirString)
	return NewResultFrom(parent)
}

func FileExtension(fileNameOrPath string) string {
	return filepath.Ext(fileNameOrPath)
}

func FileBaseName(fileNameOrPath string) string {
	return filepath.Base(fileNameOrPath)
}

func FileWithoutExtension(fileNameOrPath string) string {
	fileNameOrPath = FileBaseName(fileNameOrPath)
	var extLen = len(FileExtension(fileNameOrPath))
	var nameWithoutExtension = fileNameOrPath[0 : len(fileNameOrPath)-extLen]
	return nameWithoutExtension
}

func (dir Dir) NewIterator() Iterator[int, string] {
	var entries = NewArray[string]()
	var files, _ = os.ReadDir(dir.String())
	for _, file := range files {
		entries.Push(file.Name())
	}
	return entries.NewIterator()
}
