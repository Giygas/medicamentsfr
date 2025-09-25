// Package medicamentsparser provides functionality for downloading and parsing medicament data from external sources.
package medicamentsparser

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
)

func validateMedicament(m *entities.Medicament) error {
	if m.Cis <= 0 {
		return fmt.Errorf("invalid CIS: %d", m.Cis)
	}
	if m.Denomination == "" {
		return fmt.Errorf("missing denomination")
	}
	if m.FormePharmaceutique == "" {
		return fmt.Errorf("missing forme pharmaceutique")
	}
	// Add more checks as needed
	return nil
}

func ParseAllMedicaments() []entities.Medicament {
	// Download the neccesary files from https://base-donnees-publique.medicaments.gouv.fr/telechargement
	downloadAndParseAll()

	//Make all the json files concurrently
	var wg sync.WaitGroup
	wg.Add(5)

	conditionsChan := make(chan []entities.Condition)
	presentationsChan := make(chan []entities.Presentation)
	specialitesChan := make(chan []entities.Specialite)
	generiquesChan := make(chan []entities.Generique)
	compositionsChan := make(chan []entities.Composition)

	go func() {
		conditionsChan <- makeConditions(&wg)
	}()

	go func() {
		presentationsChan <- makePresentations(&wg)
	}()

	go func() {
		specialitesChan <- makeSpecialites(&wg)
	}()

	go func() {
		generiquesChan <- makeGeneriques(&wg)
	}()

	go func() {
		compositionsChan <- makeCompositions(&wg)
	}()

	wg.Wait()

	conditions := <-conditionsChan
	presentations := <-presentationsChan
	specialites := <-specialitesChan
	generiques := <-generiquesChan
	compositions := <-compositionsChan

	fmt.Printf("Number of conditions to process: %d\n", len(conditions))
	fmt.Printf("Number of presentations to process: %d\n", len(presentations))
	fmt.Printf("Number of generiques to process: %d\n", len(generiques))
	fmt.Printf("Number of specialites to process: %d\n", len(specialites))

	conditionsChan = nil
	presentationsChan = nil
	specialitesChan = nil
	generiquesChan = nil
	compositionsChan = nil

	var medicamentsSlice []entities.Medicament

	for _, med := range specialites {

		medicament := new(entities.Medicament)

		medicament.Cis = med.Cis
		medicament.Denomination = med.Denomination
		medicament.FormePharmaceutique = med.FormePharmaceutique
		medicament.VoiesAdministration = med.VoiesAdministration
		medicament.StatusAutorisation = med.StatusAutorisation
		medicament.TypeProcedure = med.TypeProcedure
		medicament.EtatComercialisation = med.EtatComercialisation
		medicament.DateAMM = med.DateAMM
		medicament.Titulaire = med.Titulaire
		medicament.SurveillanceRenforcee = med.SurveillanceRenforcee

		// Get all the compositions of this medicament
		for _, v := range compositions {
			if med.Cis == v.Cis {
				medicament.Composition = append(medicament.Composition, v)
			}
		}

		// Get all the generiques of this medicament
		for _, v := range generiques {
			if med.Cis == v.Cis {
				medicament.Generiques = append(medicament.Generiques, v)
			}
		}

		// Get all the presentations of this medicament
		for _, v := range presentations {
			if med.Cis == v.Cis {
				medicament.Presentation = append(medicament.Presentation, v)
			}
		}

		// Get the conditions of this medicament
		for _, v := range conditions {
			if med.Cis == v.Cis {
				medicament.Conditions = append(medicament.Conditions, v.Condition)
			}
		}

		// Validate the medicament structure
		if err := validateMedicament(medicament); err != nil {
			fmt.Printf("Skipping invalid medicament %d: %v\n", med.Cis, err)
			continue
		}

		medicamentsSlice = append(medicamentsSlice, *medicament)

	}

	fmt.Println("All medicaments parsed successfully")
	jsonMedicament, err := json.MarshalIndent(medicamentsSlice, "", "  ")
	if err != nil {
		fmt.Printf("error marshalling medicaments: %v\n", err)
		return nil
	}

	err = os.WriteFile("src/Medicaments.json", jsonMedicament, 0644)
	if err != nil {
		fmt.Printf("error writing Medicaments.json: %v\n", err)
		return nil
	}
	fmt.Println("Medicaments.json created")
	os.Stdout.Sync()

	conditions = nil
	presentations = nil
	specialites = nil
	generiques = nil
	compositions = nil

	return medicamentsSlice
}
