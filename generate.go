package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const totalFile = 3000
const contentLength = 5000

var temp = filepath.Join(os.Getenv("TEMP"), "golangccr")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Okey() {
	log.Println("start")
	start := time.Now()

	generateFiles()

	duration := time.Since(start)
	log.Println("done in", duration.Seconds(), "seconds")
}

func generateFiles() {
	os.RemoveAll(temp)
	os.MkdirAll(temp, os.ModePerm)

	for i := 0; i < totalFile; i++ {
		filename := filepath.Join(temp, fmt.Sprintf("file-%d.txt", i))
		content := RandStringAndNumber(contentLength)
		err := ioutil.WriteFile(filename, []byte(content), os.ModePerm)
		if err != nil {
			log.Println("Error writing file", filename)
		}

		if i%100 == 0 && i > 0 {
			log.Println(i, "files created")
		}
	}

	log.Printf("%d of total files created", totalFile)
}

func RandStringAndNumber(panjang int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, panjang)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
