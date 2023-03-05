package main

import (
	"fmt"
	"log"

	"github.com/Binbiubiubiu/lfs"
)

func main() {
	client, err := lfs.NewClient("baseDir")
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
