package main

import (
	"encoding/json"
	"fmt"
	"bufio"
	"os"
	"strings"
	"log"
	"strconv"
)


func makePresentations() {
	
	type Presentation struct {
		Cis 										int			`json:"Cis"`
		Cip7 										int			`json:"Cip7"`
		Libelle 								string	`json:"Libelle"`
		StatusAdministratif 		string	`json:"StatusAdministratif"`
		EtatComercialisation 		string	`json:"EtatComercialisation"`
		DateDeclaration 				string	`json:"DateDeclaration"`
		Cip13 									int			`json:"Cip13"`
		Agreement 							string	`json:"Agreement"`
		TauxRemboursement 			string	`json:"TauxRemboursement"`
		Prix 										string	`json:"Prix"`
	}
	
	tsvFile, err := os.Open("files/presentations.txt")
	if err != nil {
	}
	defer tsvFile.Close()
	
	scanner := bufio.NewScanner(tsvFile)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		
		cis, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}

		cip7, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatal(err)
		}

		cip13, err := strconv.Atoi(fields[6])
		if err != nil {
			log.Fatal(err)
		}
		
		
		//TODO look at how the prices are formatted in the medicaments database, they are a bit weird
		// Remove the points as a thousands separator
		// Change the decimal separator from comma to point
		// Verify that the prix in not null first 
		// var prix float32
		// fmt.Printf("%v",fields[9])
		// if len(fields[9]) != 0 {
		// 	comaless := strings.Replace(strings.Trim(fields[9], "."), ",", ".", -1)
		// 	// formattedPrix := strings.Replace(fields[9], ",", ".", 1)
		// 	fmt.Printf("/%v",comaless)
			
		// 	converted, err := strconv.ParseFloat(comaless, 32)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
			
		// 	prix = float32(converted)
		// } 
		
		record := Presentation {
			Cis: cis,
			Cip7: cip7,
			Libelle: fields[2],
			StatusAdministratif: fields[3],
			EtatComercialisation: fields[4],
			DateDeclaration: fields[5],
			Cip13: cip13,
			Agreement: fields[7],
			TauxRemboursement: fields[8],
			Prix: fields[9],
		}
		
		jsonRecord, err := json.Marshal(record)
		if err != nil {
			fmt.Println("error:", err)
		}
		
		fmt.Println(record)
		fmt.Println(string(jsonRecord))
		
		// var parsedPresentations []Presentation
		// fmt.Println(json.Unmarshal([]byte(jsonRecord),&parsedPresentations))
		
		fmt.Println("-----------------")
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
	// jsonData, err := json.MarshalIndent(&data, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _ = os.WriteFile("test.json", jsonData, 0644)
	
	// return jsonRecord, nil
}
