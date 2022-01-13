package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var tempPath = filepath.Join(os.Getenv("TEMP"), "golangccr")

type FileInfo struct {
	FilePath  string
	Content   []byte
	Sum       string
	IsRenamed bool
}

func main() {
	log.Println("Start")
	start := time.Now()

	chanFileContent := ReadFiles()

	chanFileSum1 := GetSum(chanFileContent)
	chanFileSum2 := GetSum(chanFileContent)
	chanFileSum3 := GetSum(chanFileContent)
	chanFileSum := MergeChannel(chanFileSum1, chanFileSum2, chanFileSum3)

	chanRename1 := Rename(chanFileSum)
	chanRename2 := Rename(chanFileSum)
	chanRename3 := Rename(chanFileSum)
	chanRename4 := Rename(chanFileSum)
	chanRename := MergeChannel(chanRename1, chanRename2, chanRename3, chanRename4)

	counterRenamed := 0
	counterTotal := 0
	for fileInfo := range chanRename {
		if fileInfo.IsRenamed {
			counterRenamed++
		}
		counterTotal++
	}

	log.Printf("%d/%d files renamed", counterRenamed, counterTotal)

	duration := time.Since(start)
	log.Println("done in", duration.Seconds(), "seconds")

}

func ReadFiles() <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		err := filepath.Walk(tempPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			chanOut <- FileInfo{
				FilePath: path,
				Content:  buf,
			}

			return nil
		})
		if err != nil {
			log.Println(err.Error())
		}
		close(chanOut)
	}()
	return chanOut
}

func GetSum(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		for fileInfo := range chanIn {
			fileInfo.Sum = fmt.Sprintf("%x", md5.Sum(fileInfo.Content))
			chanOut <- fileInfo
		}
		close(chanOut)
	}()
	return chanOut
}

func MergeChannel(chanInMany ...<-chan FileInfo) <-chan FileInfo {
	wg := new(sync.WaitGroup)
	chanOut := make(chan FileInfo)

	wg.Add(len(chanInMany))
	for _, eachChan := range chanInMany {
		go func(eachChan <-chan FileInfo) {
			for eachChanData := range eachChan {
				chanOut <- eachChanData
			}
			wg.Done()
		}(eachChan)
	}

	go func() {
		wg.Wait()
		close(chanOut)
	}()

	return chanOut
}

func Rename(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo)

	go func() {
		for fi := range chanIn {
			newPath := filepath.Join(tempPath, fmt.Sprintf("file-%s.txt", fi.Sum))
			err := os.Rename(fi.FilePath, newPath)
			fi.IsRenamed = err == nil
			chanOut <- fi
		}
		close(chanOut)
	}()
	return chanOut
}
