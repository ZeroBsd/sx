// SPDX-License-Identifier: 0BSD
package sx

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type File string

func NewFile(fullOrRelativePathAndFilename string) File {
	return File(fullOrRelativePathAndFilename)
}

func NewFileUserDir() File {
	var userHomeDir, err = os.UserHomeDir()
	ThrowIf(err != nil, ReflectFunctionName(), err.Error())
	return NewFile(userHomeDir)
}

func (f File) uri() string {
	return string(f)
}

func (f File) ReadAllBytes() []byte {
	var content, err = os.ReadFile(f.uri())
	if err != nil {
		return []byte{}
	}
	return content
}

// Reads all the text from a file, newlines are always '\n'
func (f File) ReadAllText() Optional[string] {
	var content, err = os.ReadFile(f.uri())
	if err != nil {
		return NewOptional[string]()
	}
	var text = strings.ReplaceAll(string(content), "\r\n", "\n")
	return NewOptionalFrom(text)
}

func (f File) WriteAllText(text string) {
	var err = os.WriteFile(f.uri(), []byte(text), fs.ModePerm)
	ThrowIf(err != nil, "file::WriteAllText", "Fatal error: could not write file '"+f.uri()+"'")
}

func (f File) Cd(name string) File {
	var newFileName = f.uri() + string(filepath.Separator) + name
	return NewFile(newFileName)
}

func (f File) ListDir() Array[File] {
	var d = NewArray[File]()
	var files, err = os.ReadDir(f.uri())
	if err != nil {
		return d
	}
	for _, f := range files {
		f.Name()
	}
	return d
}

func (f File) Exists() bool {
	_, err := os.Stat(f.uri())
	return err == nil
}

func (f File) IsDirectory() bool {
	r, err := os.Stat(f.uri())
	return err == nil && r.IsDir()
}

func (f File) IsSymlink() bool {
	fd, err := os.Lstat(f.uri())
	return err == nil && (fd.Mode()&os.ModeSymlink) != 0
}

func (f File) Absolute() Result[File] {
	r, err := filepath.Abs(f.uri())
	if err != nil {
		return NewResultError[File](err.Error())
	}
	return NewResultFrom(NewFile(r))
}

func (f File) Extension() File {
	return NewFile(filepath.Ext(f.uri()))
}

func (f File) Basename() File {
	return NewFile(filepath.Base(f.uri()))
}

func (f File) Canonical() File {
	return NewFile(filepath.Clean(f.uri()))
}

func (f File) RelativeTo(other File) Result[File] {
	r, err := filepath.Rel(f.uri(), other.uri())
	if err != nil {
		return NewResultError[File](err.Error())
	}
	return NewResultFrom(NewFile(r))
}
