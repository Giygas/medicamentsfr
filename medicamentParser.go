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

	// CREATE THE FULL MEDICAMENT WITH ALL HIS VARIABLES
	conditions := conditionFileToJSON()
	presentations := presentationFileToJSON()
	specialites := specialitesFileToJSON()
	generiques := generiqueFileToJSON()
	compositions := compositionFileToJSON()

	var medicamentsSlice []entities.Medicament
	
	for _, med := range specialites {
		
		PrintMemUsage()
		
		medicament := new(entities.Medicament)
		
		medicament.Cis = med.Cis
		medicament.Denomination = med.Denomination
		medicament.FormePharmaceutique = med.FormePharmaceutique
		medicament.StatusAutorisation = med.StatusAutorisation
		medicament.TypeProcedure = med.TypeProcedure
		medicament.EtatComercialisation = med.EtatComercialisation
		medicament. DateAMM = med.DateAMM
		medicament.Titulaire = med.Titulaire
		medicament.SurveillanceRenforcee = med.SurveillanceRenforcee
		
		var wg sync.WaitGroup
		
		wg.Add(4)
		// Get all the compositions of this medicament
		medicament.Composition = func (id int) []entities.Composition {
			defer wg.Done()
			var allCompos []entities.Composition
			for _, v := range (compositions) {
				if id == v.Cis {
					allCompos = append(allCompos, v)
				}
			}
			return allCompos
		}(med.Cis)
		
		// Get all the generiques of this medicament
		medicament.Generiques = func (id int) []entities.Generique {
			defer wg.Done()
			var allGeneriques []entities.Generique
			for _, v := range (generiques) {
				if id == v.Cis {
					allGeneriques = append(allGeneriques, v)
				}
			}
			return allGeneriques
		}(med.Cis)
		
		// Get all the presentations of thi medicament
		medicament.Presentation = func(id int) []entities.Presentation {
			defer wg.Done()
			var allPresentations []entities.Presentation
			for _, v := range (presentations) {
				if id == v.Cis {
					allPresentations = append(allPresentations, v)
				}
			}
			return allPresentations
		}(med.Cis)
		
		// Get the conditions of this medicament
		medicament.Conditions = func(id int) []entities.Condition {
			defer wg.Done()
			var allConditions []entities.Condition
			for _, v := range (conditions) {
				if id == v.Cis {
					allConditions = append(allConditions, v)
				}
			}
			return allConditions
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
	return medicamentsSlice
}
