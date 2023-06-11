package lfs

import (
	"fmt"
	"io"
	"log"
	"os"
	path "path/filepath"
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

func TestLocalDiskClient_UploadBuffer(t *testing.T) {
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
				return
			}
			actual, err := os.ReadFile(path.Join(dir, key))
			if err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, content, actual)
		})
	}
}

func TestLocalDiskClient_AppendBuffer(t *testing.T) {
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
				return
			}
			actual, err := os.ReadFile(path.Join(dir, key))
			if err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, tt.expected, string(actual))
		})
	}

}

func TestLocalDiskClient_Upload(t *testing.T) {
	key := "hello/upload.txt"
	if err := client.Upload(fullFilename, key); err != nil {
		t.Error(err)
		return
	}
	actual, err := os.ReadFile(path.Join(dir, key))
	if err != nil {
		t.Error(err)
		return
	}
	expected, err := os.ReadFile(fullFilename)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, expected, actual)
}

func TestLocalDiskClient_Download(t *testing.T) {
	key := "hello/download-bar.tgz"
	expected := []byte("hello bar")
	err := client.UploadBuffer(expected, key)
	if err != nil {
		t.Error(err)
		return
	}
	dest := path.Join(dir, "world")
	err = client.Download(key, dest)
	if err != nil {
		t.Error(err)
		return
	}
	actual, err := os.ReadFile(path.Join(dir, "world"))
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, expected, actual)
	_ = os.Remove(dest)
}

func TestLocalDiskClient_CreateDownloadStream(t *testing.T) {
	t.Run("should get download stream ok", func(t *testing.T) {
		key := "hello/download-bar.tgz"
		expected := []byte("hello bar")
		err := client.UploadBuffer(expected, key)
		if err != nil {
			t.Error(err)
			return
		}
		src, err := client.Open(key)
		if err != nil {
			t.Error(err)
			return
		}

		dest := path.Join(dir, "world")
		ws, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			t.Error(err)
			return
		}

		_, err = io.Copy(ws, src)
		if err != nil {
			t.Error(err)
			return
		}

		actual, err := os.ReadFile(dest)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, expected, actual)
		_ = os.Remove(dest)
	})

	t.Run("should get nil when file not exists", func(t *testing.T) {
		f, err := client.Open("hello/notexists.tgz")
		assert.True(t, os.IsNotExist(err))
		assert.True(t, f == nil)
	})
}

func TestLocalDiskClient_ReadBytes(t *testing.T) {
	t.Run("should get bytes ok", func(t *testing.T) {
		expected := []byte("hello bar")
		key := "hello/download-bar.tgz"
		err := client.UploadBuffer(expected, key)
		if err != nil {
			t.Error(err)
			return
		}
		actual, err := client.ReadFile(key)
		if err != nil {
			t.Error(err)
			return
		}
		assert.Equal(t, expected, actual)
		// _ = os.Remove(path.Join(dir, key))
	})

	t.Run("should get empty slice when file not exists", func(t *testing.T) {
		key := "hello/download-bar.tgz"
		err := os.Remove(path.Join(dir, key))
		assert.Nil(t, err, err)
		b, err := client.ReadFile(key)
		assert.True(t, os.IsNotExist(err), err)
		assert.True(t, len(b) == 0, b)
	})
}

func TestLocalDiskClient_Remove(t *testing.T) {
	content := []byte("hello bar")
	file1 := "hello/download-bar.tgz"
	file2 := "/foo/-/foo-1.3.2.txt"

	err := client.UploadBuffer(content, file1)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.UploadBuffer(content, file2)
	if err != nil {
		t.Error(err)
		return
	}

	err = client.Remove(file1)
	if err != nil {
		t.Error(err)
		return
	}
	err = client.Remove(file2)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = os.Stat(file1)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(file2)
	assert.True(t, os.IsNotExist(err))
}

func TestLocalDiskClient_List(t *testing.T) {
	err := client.Upload(fullFilename, "hello2222/upload.txt")
	if err != nil {
		t.Error(err)
		return
	}

	files, err := client.List("hello2222")
	if err != nil {
		t.Error(err)
		return
	}
	fileNames := []string{}
	for _, f := range files {
		fileNames = append(fileNames, f.Name())
	}
	assert.Equal(t, fileNames, []string{"upload.txt"})
}

func Example() {
	client, err := NewClient("baseDir")
	if err != nil {
		log.Fatal(err)
	}
	key := "test.txt"
	content := []byte("hello world")

	err = client.Upload("needUpload.txt", key)
	if err != nil {
		log.Fatal(err)
	}

	err = client.UploadBuffer(content, key)
	if err != nil {
		log.Fatal(err)
	}

	err = client.AppendBuffer(content, key)
	if err != nil {
		log.Fatal(err)
	}

	f, err := client.Open(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Open: %s", f.Name())

	bs, err := client.ReadFile(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ReadFile: %s", string(bs))

	err = client.Download(key, "download.txt")
	if err != nil {
		log.Fatal(err)
	}

	files, err := client.List("baseDir")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("List:")
	for _, f := range files {
		println(f.Name())
	}

	err = client.Remove(key)
	if err != nil {
		log.Fatal(err)
	}
}
