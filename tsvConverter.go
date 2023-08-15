package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
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
		Prix 										float32	`json:"Prix"`
	}
	
	tsvFile, err := os.Open("files/presentations.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()
	
	scanner := bufio.NewScanner(tsvFile)
	
	var jsonRecords []Presentation
	
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
		
		
		// Because the downloaded database has commas as thousands and decimal separators,
		// all the commas have to be removed except for the last one
		// If the prix is empty, 0.0 will we added in the prix section
		var prix float32
	
		if fields[9] != "" {
			
			// Count the number of commas
			numCommas := strings.Count(fields[9], ",")
			
			// If there's more than one comma, replace all but the last one
			if numCommas > 1 {
				fields[9] = strings.Replace(fields[9], ",", "", numCommas-1)
			}

			// Replace the last comma with a period
			p, err := strconv.ParseFloat(strings.Replace(fields[9], ",", ".", -1), 32)
			
			if err != nil {
				log.Fatal(err)
			}
			p = math.Trunc(p*100) / 100
			
			prix = float32(p)
		} else {
			prix = 0.0
		}
		
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
			Prix: prix,
		}
		
		jsonRecords = append(jsonRecords, record)
	}
	
	jsonData, err := json.MarshalIndent(jsonRecords, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	
	_ = os.WriteFile("src/presentations.json", jsonData, 0644)
}
