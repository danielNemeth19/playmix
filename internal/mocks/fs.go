package mocks

import (
	"errors"
	"io/fs"
	"time"
)

type FakeReadError struct {
	message string
}

func (f FakeReadError) Error() string {
    return f.message
}

type FakeFile struct{}

func (f FakeFile) Stat() (fs.FileInfo, error) {
	error := fs.PathError{Op: "faked stat", Path: "/dummy/file", Err: fs.ErrInvalid}
	return nil, &error
}

func (f FakeFile) Read([]byte) (int, error) {
    error := FakeReadError{message: "Faked read error"}
	return 10, &error
}

func (f FakeFile) Close() error {
	return nil
}

type FakeSys struct{}

func (f FakeSys) Open(name string) (fs.File, error) {
	return FakeFile{}, nil
}

type FakeFileInfo struct {
	name     string
	size     int64
	mode     fs.FileMode
	modeTime time.Time
	isDir    bool
}

func (m FakeFileInfo) Name() string {
	return m.name
}

func (m FakeFileInfo) Size() int64 {
	return m.size
}

func (m FakeFileInfo) Mode() fs.FileMode {
	return m.mode
}

func (m FakeFileInfo) ModTime() time.Time {
	return m.modeTime
}

func (m FakeFileInfo) IsDir() bool {
	return m.isDir
}

func (m FakeFileInfo) Sys() any {
	return nil
}

type FakeDirEntry struct {
	name     string
	fileInfo fs.FileInfo
	isDir    bool
}

func (m FakeDirEntry) Name() string {
	return m.name
}

func (m FakeDirEntry) IsDir() bool {
	return m.isDir
}

func (m FakeDirEntry) Type() fs.FileMode {
	return 0755
}

func (m FakeDirEntry) Info() (fs.FileInfo, error) {
	return m.fileInfo, nil
}

func CreateFakeDirEntry(name string, isDir bool, modTime time.Time) fs.DirEntry {
	ff := FakeFileInfo{
		name:     name,
		modeTime: modTime,
		isDir:    isDir,
	}
	fd := FakeDirEntry{
		name:     ff.Name(),
		fileInfo: ff,
		isDir:    isDir,
	}
	return fd
}

type FakeWriter struct{}

func (fw FakeWriter) Write(p []byte) (n int, err error) {
	error := errors.New("fake error for io.Writer Write function")
	return 0, error
}
