package mocks

import (
	"io/fs"
	"time"
)

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
    isDir bool
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
        isDir: isDir,
	}
	fd := FakeDirEntry{
		name:     ff.Name(),
		fileInfo: ff,
        isDir: isDir,
	}
	return fd
}
