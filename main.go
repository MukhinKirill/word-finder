package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 2 {
		fmt.Println("Не указанны все обязательные аргументы при запуске утилиты.")
		fmt.Println("Пример: word-finder.exe Path(обязательный) word(обязательный) charset(дополнительный)")
		os.Exit(1)
	}
	folderPath := argsWithoutProg[0]
	wordForFind := strings.ToLower(argsWithoutProg[1])
	fileCharset := "utf-8"
	if len(argsWithoutProg) > 2 {
		fileCharset = strings.ToLower(argsWithoutProg[2])
	}

	if folderPath[len(folderPath)-1] != '\\' {
		folderPath = folderPath + "\\"
	}

	fmt.Println(argsWithoutProg)
	fmt.Println("Search word:", wordForFind)

	files := getAllfiles(folderPath)
	start := time.Now()

	var ch = make(chan int, len(files))
	for _, f := range files {

		go func(fi os.FileInfo, ch chan<- int) {
			if fi.Mode().IsRegular() {
				filePath := folderPath + fi.Name()
				wordCount := findWordInFile(filePath, fileCharset, wordForFind)
				fmt.Println(fi.Name(), wordCount)
				ch <- wordCount
			}
		}(f, ch)
	}
	wordCounter := 0
	for i := 0; i < len(files); i++ {
		wordCounter += <-ch
	}
	close(ch)
	elapsed := time.Since(start)
	fmt.Printf("\nКоличество найденных слов во всех файлах: %d", wordCounter)
	fmt.Printf("\nЗатраченное время %s", elapsed)
}

func getAllfiles(path string) []os.FileInfo {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func findWordInFile(filePath string, fileCharset string, searchedWord string) int {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
	text := string(data)

	reader := strings.NewReader(text)

	contentType := fmt.Sprint("text/html; charset = ", fileCharset)
	utf8Reader, err := charset.NewReader(reader, contentType)
	if err != nil {
		fmt.Println("Encoding error:", err)
		return -1
	}

	utf8Bytes, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		fmt.Println("IO error:", err)
		return -1
	}

	utf8Text := string(utf8Bytes)
	wordCounter := 0
	words := strings.Fields(utf8Text)
	for _, word := range words {
		if strings.ToLower(word) == searchedWord {
			wordCounter++
		}
	}
	return wordCounter
}
