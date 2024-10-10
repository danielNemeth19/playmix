package mocks

import (
	"io/fs"
	"time"
)

type FakeFileInfo struct {
	FileName     string
	SizeField    int64
	ModeField    fs.FileMode
	ModTimeField time.Time
}

func (m FakeFileInfo) Name() string {
	return m.FileName
}

func (m FakeFileInfo) Size() int64 {
	return m.SizeField
}

func (m FakeFileInfo) Mode() fs.FileMode {
	return m.ModeField
}

func (m FakeFileInfo) ModTime() time.Time {
	return m.ModTimeField
}

func (m FakeFileInfo) IsDir() bool {
	return false
}

func (m FakeFileInfo) Sys() any {
	return nil
}

type FakeDirEntry struct {
	NameField string
	FileInfo  fs.FileInfo
}

func (m FakeDirEntry) Name() string {
	return m.NameField
}

func (m FakeDirEntry) IsDir() bool {
	return false
}

func (m FakeDirEntry) Type() fs.FileMode {
	return 0755
}

func (m FakeDirEntry) Info() (fs.FileInfo, error) {
	return m.FileInfo, nil
}

func CreateFakeDirEntry(name string, modTime time.Time) fs.DirEntry {
	ff := FakeFileInfo{
		FileName:     name,
		ModTimeField: modTime,
	}
	fd := FakeDirEntry{
		NameField: ff.Name(),
		FileInfo:  ff,
	}
	return fd
}
