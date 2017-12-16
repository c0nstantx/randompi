package services

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/h2non/filetype.v1"
)

type videoFile struct {
	Hash string
	Name string
	Path string
}

func VideoList(path string) map[string]videoFile {
	return readFiles(path)
}

func readFiles(path string) map[string]videoFile {
	fileList := make(map[string]videoFile)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fileName := file.Name()
		filePath := filepath.Join(path, fileName)
		if file.IsDir() {
			folderFiles := readFiles(filePath)
			for _, folderFile := range folderFiles {
				fileList[folderFile.Hash] = folderFile
			}
		} else {
			f := buildVideoFile(fileName, filePath)
			if len(f.Name) > 0 {
				fileList[f.Hash] = f
			}
		}
	}

	return fileList
}

func buildVideoFile(name string, path string) videoFile {
	vFile := videoFile{}
	file, _ := os.Open(path)
	head := make([]byte, 261)
	file.Read(head)
	if filetype.IsVideo(head) {
		hashData := []byte(path)
		hashBytes := md5.Sum(hashData)
		hash := fmt.Sprintf("%x", hashBytes)
		vFile.Hash = hash
		vFile.Name = name
		vFile.Path = path
	}
	file.Close()

	return vFile
}
