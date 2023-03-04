package lfs

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client *LocalDiskClient
var dir string
var fullFilename string

func init() {
	cwd, _ := os.Getwd()
	dir = path.Join(cwd, "fixtures")
	_ = os.RemoveAll(dir)
	client, _ = NewClient(dir)
	_, fullFilename, _, _ = runtime.Caller(0)
}

func TestNewClient(t *testing.T) {

	t.Run("should throw error with empty path", func(t *testing.T) {
		_, err := NewClient("")
		assert.EqualError(t, err, "need present dir")
	})

	t.Run("should create client ok", func(t *testing.T) {
		client, _ = NewClient(dir)
		assert.Equal(t, dir, client.dir)
	})
}

func Test_UploadBuffer(t *testing.T) {
	content := []byte("hello")
	tests := []string{
		"hello/bar.tgz",
		"/a/b/c/d/e/f/g.txt",
		"/foo/-/foo-1.3.2.txt",
	}
	for i, key := range tests {
		t.Run("case "+strconv.Itoa(i), func(t *testing.T) {
			if err := client.UploadBuffer(content, key); err != nil {
				t.Error(err)
			}
			actual, err := os.ReadFile(path.Join(dir, key))
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, content, actual)
		})
	}
}

func Test_AppendBuffer(t *testing.T) {
	key := "hello/bar.txt"
	tests := []struct {
		append   string
		expected string
	}{
		{"hello", "hello"},
		{" world", "hello world"},
		{"\nagain", "hello world\nagain"},
	}
	for i, tt := range tests {
		t.Run("case "+strconv.Itoa(i), func(t *testing.T) {
			if err := client.AppendBuffer([]byte(tt.append), key); err != nil {
				t.Error(err)
			}
			actual, err := os.ReadFile(path.Join(dir, key))
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.expected, string(actual))
		})
	}

}

func Test_Upload(t *testing.T) {
	key := "hello/upload.go"
	if err := client.Upload(fullFilename, key); err != nil {
		t.Error(err)
	}
	actual, err := os.ReadFile(path.Join(dir, key))
	if err != nil {
		t.Error(err)
	}
	expected, err := os.ReadFile(fullFilename)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, actual)
}

func Test_Download(t *testing.T) {
	key := "hello/download-bar.tgz"
	expected := []byte("hello bar")
	err := client.UploadBuffer(expected, key)
	if err != nil {
		t.Error(err)
	}
	dest := path.Join(dir, "world")
	err = client.Download(key, dest)
	if err != nil {
		t.Error(err)
	}
	actual, err := os.ReadFile(path.Join(dir, "world"))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expected, actual)
	_ = os.Remove(dest)
}

// func Test_CreateDownloadStream(t *testing.T) {
// 	key := "hello/download-bar.tgz"
// 	expected := []byte("hello bar")
// 	err := client.UploadBuffer(expected, key)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	dest := path.Join(dir, "world")
// 	writeStream, err := client.Open(dest)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	actual, err := os.ReadFile(path.Join(dir, "world"))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	assert.Equal(t, expected, actual)
// 	_ = os.Remove(dest)
// }
