package filestore

import (
	"fmt"
	"sync"
)

type FileInfo struct {
	Name    string
	Content string
	Hash    string
}

type FileStore struct {
	Files map[string]FileInfo
	mu    sync.RWMutex
}

func NewFileStore() *FileStore {
	return &FileStore{
		Files: make(map[string]FileInfo),
	}
}

func (fs *FileStore) Add(filename, content string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if _, exists := fs.Files[filename]; exists {
		return fmt.Errorf("file already exists")
	}
	fs.Files[filename] = FileInfo{Name: filename, Content: content}
	return nil
}

func (fs *FileStore) Remove(filename string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	if _, exists := fs.Files[filename]; !exists {
		return fmt.Errorf("file not found")
	}
	delete(fs.Files, filename)
	return nil
}

func (fs *FileStore) Update(filename, content string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.Files[filename] = FileInfo{Name: filename, Content: content}
	return nil
}

func (fs *FileStore) List() []string {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	var filenames []string
	for name := range fs.Files {
		filenames = append(filenames, name)
	}
	return filenames
}

func (fs *FileStore) Get(filename string) (FileInfo, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	file, ok := fs.Files[filename]
	return file, ok
}
