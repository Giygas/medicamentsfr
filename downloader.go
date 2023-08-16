package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/text/encoding/charmap"
)

func downloadAndParseFile(filepath string, url string) error {
	
	filepath = "files/" + filepath + ".txt"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	outFile, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	
	defer outFile.Close()
	reader := charmap.Windows1252.NewDecoder().Reader(response.Body)
	scanner := bufio.NewScanner(reader)
	
	for scanner.Scan() {
		_, err = fmt.Fprintln(outFile, scanner.Text())
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Downloaded: " + filepath)
	return nil
}

//Download all files concurrently
//Params:
// 1. files map[string]string - A map containing the name of the file and the url
func downloadAndParseAll(files map[string]string) error {
	
	//Create the files directory if it doesn't exists
	filePath := filepath.Join(".", "files")
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic( err)
	}
	filePath = filepath.Join(".", "src")
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic( err)
	}
	
	var wg sync.WaitGroup
	
	for fileName, url := range(files) {
		wg.Add(1)
		
		go func(file string, url string) {
			defer wg.Done()
			downloadAndParseFile(file, url)
		} (fileName, url)
		
	}
	wg.Wait()

	return nil
}
