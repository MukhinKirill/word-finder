package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
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

	fmt.Println(argsWithoutProg)
	fmt.Println("Search word:", wordForFind)

	files := getAllfiles(folderPath)
	commonWordCounter := 0
	start := time.Now()

	var wg sync.WaitGroup
	for _, f := range files {
		wg.Add(1)
		go func(fi os.FileInfo) {
			defer wg.Done()
			if fi.Mode().IsRegular() {
				if folderPath[len(folderPath)-1] != '\\' {
					folderPath = folderPath + "\\"
				}
				filePath := folderPath + fi.Name()
				wordCount := findWordInFile(filePath, fileCharset, wordForFind)
				fmt.Println(fi.Name(), wordCount)
				commonWordCounter += wordCount
			}
		}(f)
	}
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("\nКоличество найденных слов во всех файлах: %d", commonWordCounter)
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
	utf8, err := charset.NewReader(reader, contentType)
	if err != nil {
		fmt.Println("Encoding error:", err)
		return -1
	}

	utf8Bytes, err := ioutil.ReadAll(utf8)
	if err != nil {
		fmt.Println("IO error:", err)
		return -1
	}

	textUtf8 := string(utf8Bytes)
	wordCounter := 0
	words := strings.Fields(textUtf8)
	for _, word := range words {
		if strings.ToLower(word) == searchedWord {
			wordCounter++
		}
	}
	return wordCounter
}
