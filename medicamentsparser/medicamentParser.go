package medicamentsparser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
)

func ParseAllMedicaments() []entities.Medicament {
	// Download the neccesary files from https://base-donnees-publique.medicaments.gouv.fr/telechargement.php
	if err := downloadAndParseAll(); err != nil {
		log.Fatalf("Failed to download files: %v", err)
	}

	//Make all the json files concurrently
	var wg sync.WaitGroup
	wg.Add(5)

	type result struct {
		data interface{}
		err  error
	}

	conditionsChan := make(chan result)
	presentationsChan := make(chan result)
	specialitesChan := make(chan result)
	generiquesChan := make(chan result)
	compositionsChan := make(chan result)

	go func() {
		defer wg.Done()
		data, err := makeConditions(nil)
		conditionsChan <- result{data, err}
	}()

	go func() {
		defer wg.Done()
		data, err := makePresentations(nil)
		presentationsChan <- result{data, err}
	}()

	go func() {
		defer wg.Done()
		data, err := makeSpecialites(nil)
		specialitesChan <- result{data, err}
	}()

	go func() {
		defer wg.Done()
		data, err := makeGeneriques(nil)
		generiquesChan <- result{data, err}
	}()

	go func() {
		defer wg.Done()
		data, err := makeCompositions(nil)
		compositionsChan <- result{data, err}
	}()

	wg.Wait()

	conditionsRes := <-conditionsChan
	if conditionsRes.err != nil {
		log.Fatalf("Failed to make conditions: %v", conditionsRes.err)
	}
	conditions := conditionsRes.data.([]entities.Condition)

	presentationsRes := <-presentationsChan
	if presentationsRes.err != nil {
		log.Fatalf("Failed to make presentations: %v", presentationsRes.err)
	}
	presentations := presentationsRes.data.([]entities.Presentation)

	specialitesRes := <-specialitesChan
	if specialitesRes.err != nil {
		log.Fatalf("Failed to make specialites: %v", specialitesRes.err)
	}
	specialites := specialitesRes.data.([]entities.Specialite)

	generiquesRes := <-generiquesChan
	if generiquesRes.err != nil {
		log.Fatalf("Failed to make generiques: %v", generiquesRes.err)
	}
	generiques := generiquesRes.data.([]entities.Generique)

	compositionsRes := <-compositionsChan
	if compositionsRes.err != nil {
		log.Fatalf("Failed to make compositions: %v", compositionsRes.err)
	}
	compositions := compositionsRes.data.([]entities.Composition)

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

		var wg sync.WaitGroup

		wg.Add(4)
		// Get all the compositions of this medicament
		go func(id int) {
			defer wg.Done()
			for _, v := range compositions {
				if id == v.Cis {
					medicament.Composition = append(medicament.Composition, v)
				}
			}
		}(med.Cis)

		// Get all the generiques of this medicament
		go func(id int) {
			defer wg.Done()
			for _, v := range generiques {
				if id == v.Cis {
					medicament.Generiques = append(medicament.Generiques, v)
				}
			}
		}(med.Cis)

		// Get all the presentations of this medicament
		go func(id int) {
			defer wg.Done()
			for _, v := range presentations {
				if id == v.Cis {
					medicament.Presentation = append(medicament.Presentation, v)
				}
			}
		}(med.Cis)

		// Get the conditions of this medicament
		go func(id int) {
			defer wg.Done()
			for _, v := range conditions {
				if id == v.Cis {
					medicament.Conditions = append(medicament.Conditions, v.Condition)
				}
			}
		}(med.Cis)

		wg.Wait()
		medicamentsSlice = append(medicamentsSlice, *medicament)

	}

	jsonMedicament, err := json.MarshalIndent(medicamentsSlice, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}

	_ = os.WriteFile("src/Medicaments.json", jsonMedicament, 0644)
	log.Println("Medicaments.json created")

	conditions = nil
	presentations = nil
	specialites = nil
	generiques = nil
	compositions = nil

	return medicamentsSlice
}
