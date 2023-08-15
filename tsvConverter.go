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
	"sync"
)


func makePresentations(wg *sync.WaitGroup) {
	defer wg.Done()
	type Presentation struct {
		Cis 										int			`json:"cis"`
		Cip7 										int			`json:"cip7"`
		Libelle 								string	`json:"libelle"`
		StatusAdministratif 		string	`json:"statusAdministratif"`
		EtatComercialisation 		string	`json:"etatComercialisation"`
		DateDeclaration 				string	`json:"dateDeclaration"`
		Cip13 									int			`json:"cip13"`
		Agreement 							string	`json:"agreement"`
		TauxRemboursement 			string	`json:"tauxRemboursement"`
		Prix 										float32	`json:"prix"`
	}
	
	tsvFile, err := os.Open("files/Presentations.txt")
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
	
	_ = os.WriteFile("src/Presentations.json", jsonData, 0644)
	log.Println("Presentations.json created")
}

func makeGeneriques(wg *sync.WaitGroup) {
	defer wg.Done()
	type Generique struct {
		Cis 										int			`json:"cis"`
		Group 									int			`json:"group"`
		Libelle 								string	`json:"libelle"`
	}
	
	tsvFile, err := os.Open("files/Generiques.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()
	
	scanner := bufio.NewScanner(tsvFile)
	
	var jsonRecords []Generique
	
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		
		cis, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Fatal(err)
		}
		
		group, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}
		
		record := Generique {
			Cis: cis,
			Group: group,
			Libelle: fields[1],
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
	
	_ = os.WriteFile("src/Generiques.json", jsonData, 0644)
	log.Println("Generiques.json created")
}

func makeCompositions(wg *sync.WaitGroup) {
	defer wg.Done()
	type Composition struct {
		Cis 										int			`json:"cis"`
		ElementParmaceutique 		string	`json:"elementPharmaceutique"`
		CodeSubstance 					int			`json:"codeSubstance"`
		VoiesAdministration 	string	`json:"VoiesAdministration"`
		Dosage 									string	`json:"dosage"`
		ReferenceDosage 				string	`json:"referenceDosage"`
		NatureComposant 				string	`json:"natureComposant"`
	}
	
	tsvFile, err := os.Open("files/Compositions.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()
	
	scanner := bufio.NewScanner(tsvFile)
	
	var jsonRecords []Composition
	
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		
		cis, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}

		codeS, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Fatal(err)
		}
		
		record := Composition {
			Cis: cis,
			ElementParmaceutique: fields[1],
			CodeSubstance: codeS,
			VoiesAdministration: fields[3],
			Dosage: fields[4],
			ReferenceDosage: fields[5],
			NatureComposant: fields[6],
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
	
	_ = os.WriteFile("src/Compositions.json", jsonData, 0644)
	log.Println("Compositions.json created")
}

func makeSpecialites(wg *sync.WaitGroup) {
	defer wg.Done()
	type Specialite struct {
		Cis 										int			`json:"cis"`
		Denomination 						string	`json:"elementPharmaceutique"`
		FormePharmaceutique 		string	`json:"formePharmaceutique"`
		VoiesAdministration 		string	`json:"voiesAdministration"`
		StatusAutorisation 			string	`json:"statusAutorisation"`
		TypeProcedure	 					string	`json:"typeProcedure"`
		EtatComercialisation 		string	`json:"etatComercialisation"`
		DateAMM 								string	`json:"dateAMM"`
		Titulaire 							string	`json:"titulaire"`
		SurveillanceRenforcee		string	`json:"surveillanceRenforce"`
	}
	
	tsvFile, err := os.Open("files/Specialites.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()
	
	scanner := bufio.NewScanner(tsvFile)
	
	var jsonRecords []Specialite
	
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		
		cis, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}
		
		record := Specialite {
			Cis: cis,
			Denomination: fields[1],
			FormePharmaceutique: fields[2],
			VoiesAdministration: fields[3],
			StatusAutorisation: fields[4],
			TypeProcedure: fields[5],
			EtatComercialisation: fields[6],
			DateAMM: fields[7],
			Titulaire: fields[10],
			SurveillanceRenforcee: fields[11],
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
	
	_ = os.WriteFile("src/Specialites.json", jsonData, 0644)
	log.Println("Specialites.json created")
}

func makeConditions(wg *sync.WaitGroup) {
	defer wg.Done()
	type Condition struct {
		Cis 										int			`json:"cis"`
		Condition 							string	`json:"condition"`
	}
	
	tsvFile, err := os.Open("files/Conditions.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()
	
	scanner := bufio.NewScanner(tsvFile)
	
	var jsonRecords []Condition
	
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		
		cis, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}
		
		record := Condition {
			Cis: cis,
			Condition: fields[1],
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
	
	_ = os.WriteFile("src/Conditions.json", jsonData, 0644)
	log.Println("Conditions.json created")
}
