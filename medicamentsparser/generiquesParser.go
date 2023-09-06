package medicamentsparser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
)

var generiquesList []entities.GeneriqueList

func GeneriquesParser(medicaments *[]entities.Medicament) []entities.GeneriqueList {

	// Create a map of all the medicaments cis to reduce algorithm complexity
	medicamentMap := make(map[int]*entities.Medicament)
	for i := range *medicaments {
		medicamentMap[(*medicaments)[i].Cis] = &(*medicaments)[i]
	}

	// allGeneriques: []Generique
	allGeneriques := makeGeneriques(nil)
	// Create a map of all the generiques to reduce algorithm complexity
	generiquesMap := make(map[int]entities.Generique)
	for i := range allGeneriques {
		generiquesMap[allGeneriques[i].Group] = allGeneriques[i]
	}

	// generiques file: [groupid]:[]cis of medicaments in the same group
	generiquesFile := generiqueFileToJSON()

	var generiques []entities.GeneriqueList

	for i, v := range generiquesFile {

		// Convert the string index to integer
		groupInt, err := strconv.Atoi(i)
		if err != nil {
			log.Println("An error ocurred converting the generiques group to integer", err)
		}

		current := entities.GeneriqueList{
			GroupId:     groupInt,
			Libelle:     generiquesMap[groupInt].Libelle,
			Medicaments: getMedicamentsInArray(v, medicaments, medicamentMap),
		}

		generiques = append(generiques, current)
	}

	marshalledGeneriques, err := json.MarshalIndent(generiques, "", " ")
	if err != nil {
		log.Println("An error has occurred when marshalling generiques", err)
	}
	_ = os.WriteFile("src/GeneriquesFull.json", marshalledGeneriques, 0644)
	fmt.Println("GeneriquesFull.json created")
	return generiques
}

func getMedicamentsInArray(medicamentsIds []int, medicaments *[]entities.Medicament, medicamentMap map[int]*entities.Medicament) []entities.GeneriqueMedicament {
	var medicamentsArray []entities.GeneriqueMedicament

	for _, v := range medicamentsIds {
		if medicament, ok := medicamentMap[v]; ok {
			generiqueMed := entities.GeneriqueMedicament{
				Cis:                 (medicament.Cis),
				Denomination:        (medicament.Denomination),
				FormePharmaceutique: (medicament.FormePharmaceutique),
			}
			medicamentsArray = append(medicamentsArray, generiqueMed)
		}
	}

	return medicamentsArray
}
