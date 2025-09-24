package medicamentsparser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
)

var medsType map[int]string

func GeneriquesParser(medicaments *[]entities.Medicament, mMap *map[int]entities.Medicament) ([]entities.GeneriqueList, map[int]entities.Generique) {
	var err error

	fmt.Println("trying to parse generiques")
	// allGeneriques: []Generique
	allGeneriques, err := makeGeneriques(nil)
	if err != nil {
		log.Fatalf("Failed to make generiques: %v", err)
	}
	// Create a map of all the generiques to reduce algorithm complexity
	generiquesMap := make(map[int]entities.Generique)
	for i := range allGeneriques {
		generiquesMap[allGeneriques[i].Group] = allGeneriques[i]
	}

	// generiques file: [groupid]:[]cis of medicaments in the same group
	generiquesFile, err := generiqueFileToJSON()
<<<<<<< HEAD
	if err != nil {
		log.Fatalf("Failed to read generiques file: %v", err)
	}
=======
	//TODO: handle the error here
>>>>>>> working-one

	// The medsType is a map where the key are the medicament cis and the value is the
	// type of generique
	medsType, err = createMedicamentGeneriqueType()
	if err != nil {
		log.Fatalf("Failed to create medicament generique type: %v", err)
	}

	var generiques []entities.GeneriqueList

	for i, v := range generiquesFile {

		// Convert the string index to integer
		groupInt, convErr := strconv.Atoi(i)
		if err != nil {
			log.Println("An error ocurred converting the generiques group to integer", convErr)
<<<<<<< HEAD
			continue
=======
>>>>>>> working-one
		}

		current := entities.GeneriqueList{
			GroupId:     groupInt,
			Libelle:     generiquesMap[groupInt].Libelle,
			Medicaments: getMedicamentsInArray(v, mMap),
		}

		generiques = append(generiques, current)
	}

	// Write debug file
	marshalledGeneriques, err := json.MarshalIndent(generiques, "", " ")
	if err != nil {
		log.Printf("Error marshalling generiques: %v", err)
	} else {
		if writeErr := os.WriteFile("src/GeneriquesFull.json", marshalledGeneriques, 0644); writeErr != nil {
			log.Printf("Error writing GeneriquesFull.json: %v", writeErr)
		} else {
			log.Println("GeneriquesFull.json created")
		}
	}

	log.Println("Generiques parsing completed")
	return generiques, generiquesMap
}

func createGeneriqueComposition(medicamentComposition *[]entities.Composition) []entities.GeneriqueComposition {
	var compositions []entities.GeneriqueComposition
	for _, v := range *medicamentComposition {
		compo := entities.GeneriqueComposition{
			ElementParmaceutique:  v.ElementParmaceutique,
			DenominationSubstance: v.DenominationSubstance,
			Dosage:                v.Dosage,
		}
		compositions = append(compositions, compo)
	}
	return compositions
}

func getMedicamentsInArray(medicamentsIds []int, medicamentMap *map[int]entities.Medicament) []entities.GeneriqueMedicament {
	var medicamentsArray []entities.GeneriqueMedicament

	for _, v := range medicamentsIds {
		if medicament, ok := (*medicamentMap)[v]; ok {
			generiqueComposition := createGeneriqueComposition(&medicament.Composition)
			generiqueMed := entities.GeneriqueMedicament{
				Cis:                 medicament.Cis,
				Denomination:        medicament.Denomination,
				FormePharmaceutique: medicament.FormePharmaceutique,
				Type:                medsType[medicament.Cis],
				Composition:         generiqueComposition,
			}
			medicamentsArray = append(medicamentsArray, generiqueMed)
		}
	}

	return medicamentsArray
}
