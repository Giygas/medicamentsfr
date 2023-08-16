package main

import (
	"encoding/json"
	"fmt"
	"log"
	"medicamentsfr/entities"
	"os"
	"sync"
)

func parseAllMedicaments() []entities.Medicament{
	
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
		medicament. DateAMM = med.DateAMM
		medicament.Titulaire = med.Titulaire
		medicament.SurveillanceRenforcee = med.SurveillanceRenforcee
		
		var wg sync.WaitGroup
		
		wg.Add(4)
		// Get all the compositions of this medicament
		go func (id int) {
			defer wg.Done()
			for _, v := range (compositions) {
				if id == v.Cis {
					medicament.Composition = append(medicament.Composition, v)
				}
			}
		}(med.Cis)
		
		// Get all the generiques of this medicament
		go func (id int){
			defer wg.Done()
			for _, v := range (generiques) {
				if id == v.Cis {
					medicament.Generiques = append(medicament.Generiques, v)
				}
			}
		}(med.Cis)
		
		// Get all the presentations of thi medicament
		go func(id int) {
			defer wg.Done()
			for _, v := range (presentations) {
				if id == v.Cis {
					medicament.Presentation = append(medicament.Presentation, v)
				}
			}
		}(med.Cis)
		
		// Get the conditions of this medicament
		go func(id int) {
			defer wg.Done()
			for _, v := range (conditions) {
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
