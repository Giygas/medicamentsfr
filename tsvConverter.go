package main

import (
	"encoding/json"
	"fmt"
	"bufio"
	"os"
	"strings"
	"log"
)

func splitTSVFile(fileName string) ([]map[int]string, error) {
	
	tsvFile, err := os.Open("files/" + fileName + ".txt")
	if err != nil {
		return nil, err
	}
	defer tsvFile.Close()
	
	scanner := bufio.NewScanner(tsvFile)

	data := make([]map[int]string, 0)
	
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		
		record := make(map[int]string)
		for i, value := range(fields){
			record[i] = value
		}
		
		data = append(data, record)
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
	//TODO to this in the main program
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))
	return data,nil
}
