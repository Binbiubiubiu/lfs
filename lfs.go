package lfs

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

type LocalDiskClient struct {
	dir string
}

func (c *LocalDiskClient) Upload(filepath string, key string) (err error) {
	destpath := c.getPath(key)
	err = c.ensureDirExists(destpath)
	if err != nil {
		return
	}
	content, err := os.ReadFile(filepath)
	if err != nil {
		return
	}
	return os.WriteFile(destpath, content, 0666)
}

func (c *LocalDiskClient) UploadBuffer(content []byte, key string) (err error) {
	destpath := c.getPath(key)
	err = c.ensureDirExists(destpath)
	if err != nil {
		return
	}
	return os.WriteFile(destpath, content, 0777)
}

func (c *LocalDiskClient) AppendBuffer(content []byte, key string) (err error) {
	destpath := c.getPath(key)
	err = c.ensureDirExists(destpath)
	if err != nil {
		return
	}
	file, err := os.OpenFile(destpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	_, err = file.Write(content)
	if err != nil {
		return
	}
	return file.Close()
}

func (c *LocalDiskClient) Open(key string) (*os.File, error) {
	filepath := c.getPath(key)
	return os.OpenFile(filepath, os.O_RDONLY, 0777)
}

func (c *LocalDiskClient) ReadFile(key string) ([]byte, error) {
	filepath := c.getPath(key)
	return os.ReadFile(filepath)
}

func (c *LocalDiskClient) Download(key string, savePath string) (err error) {
	filepath := c.getPath(key)
	content, err := os.ReadFile(filepath)
	if err != nil {
		return
	}
	return os.WriteFile(savePath, content, 0666)
}

func (c *LocalDiskClient) Remove(key string) (err error) {
	filepath := c.getPath(key)
	return os.Remove(filepath)
}

func (c *LocalDiskClient) List(prefix string) ([]fs.DirEntry, error) {
	destpath := c.getPath(prefix)
	return os.ReadDir(destpath)
}

func (c *LocalDiskClient) ensureDirExists(filepath string) (err error) {
	err = os.MkdirAll(path.Dir(filepath), 0777)
	if err != nil && os.IsExist(err) {
		return nil
	}
	return
}

func (c *LocalDiskClient) getPath(key string) string {
	return path.Join(c.dir, key)
}

func NewClient(dir string) (*LocalDiskClient, error) {
	if dir == "" {
		return nil, errors.New("need present dir")
	}

	client := LocalDiskClient{
		dir: dir,
	}
	return &client, nil
}
