package medicamentsparser

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
)

func makePresentations(wg *sync.WaitGroup) []entities.Presentation {
	defer wg.Done()

	tsvFile, err := os.Open("files/Presentations.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()

	scanner := bufio.NewScanner(tsvFile)

	var jsonRecords []entities.Presentation

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

		record := entities.Presentation{
			Cis:                  cis,
			Cip7:                 cip7,
			Libelle:              fields[2],
			StatusAdministratif:  fields[3],
			EtatComercialisation: fields[4],
			DateDeclaration:      fields[5],
			Cip13:                cip13,
			Agreement:            fields[7],
			TauxRemboursement:    fields[8],
			Prix:                 prix,
		}

		jsonRecords = append(jsonRecords, record)
	}

	log.Println("Presentations done")
	return jsonRecords
}

func makeGeneriques(wg *sync.WaitGroup) []entities.Generique {
	defer wg.Done()

	tsvFile, err := os.Open("files/Generiques.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()

	scanner := bufio.NewScanner(tsvFile)

	var jsonRecords []entities.Generique

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

		record := entities.Generique{
			Cis:     cis,
			Group:   group,
			Libelle: fields[1],
		}

		jsonRecords = append(jsonRecords, record)
	}

	log.Println("Generiques done")
	return jsonRecords
}

func makeCompositions(wg *sync.WaitGroup) []entities.Composition {
	defer wg.Done()

	tsvFile, err := os.Open("files/Compositions.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()

	scanner := bufio.NewScanner(tsvFile)

	var jsonRecords []entities.Composition

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

		record := entities.Composition{
			Cis:                   cis,
			ElementParmaceutique:  fields[1],
			CodeSubstance:         codeS,
			DenominationSubstance: fields[3],
			Dosage:                fields[4],
			ReferenceDosage:       fields[5],
			NatureComposant:       fields[6],
		}

		jsonRecords = append(jsonRecords, record)
	}

	log.Println("Compositions done")
	return jsonRecords
}

func makeSpecialites(wg *sync.WaitGroup) []entities.Specialite {
	defer wg.Done()

	tsvFile, err := os.Open("files/Specialites.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()

	scanner := bufio.NewScanner(tsvFile)

	var jsonRecords []entities.Specialite

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		cis, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}

		record := entities.Specialite{
			Cis:                   cis,
			Denomination:          fields[1],
			FormePharmaceutique:   fields[2],
			VoiesAdministration:   strings.Split(fields[3], ";"),
			StatusAutorisation:    fields[4],
			TypeProcedure:         fields[5],
			EtatComercialisation:  fields[6],
			DateAMM:               fields[7],
			Titulaire:             strings.TrimLeft(fields[10], " "),
			SurveillanceRenforcee: fields[11],
		}

		jsonRecords = append(jsonRecords, record)
	}

	log.Println("Specialites done")
	return jsonRecords
}

func makeConditions(wg *sync.WaitGroup) []entities.Condition {
	defer wg.Done()

	tsvFile, err := os.Open("files/Conditions.txt")
	if err != nil {
		log.Fatal("Error opening file", err)
	}
	defer tsvFile.Close()

	scanner := bufio.NewScanner(tsvFile)

	var jsonRecords []entities.Condition

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")

		cis, err := strconv.Atoi(fields[0])
		if err != nil {
			log.Fatal(err)
		}

		record := entities.Condition{
			Cis:       cis,
			Condition: fields[1],
		}

		jsonRecords = append(jsonRecords, record)
	}

	log.Println("Conditions done")
	return jsonRecords
}
