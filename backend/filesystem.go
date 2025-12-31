package main

import "os"

type FileSystemOps interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Create(name string) (*os.File, error)
	Open(name string) (*os.File, error)
	Remove(name string) error
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
}

type osFileSystem struct{}

func (fs *osFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (fs *osFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (fs *osFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *osFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (fs *osFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fs *osFileSystem) Open(name string) (*os.File, error) {
	return os.Open(name)
}

func (fs *osFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (fs *osFileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}
