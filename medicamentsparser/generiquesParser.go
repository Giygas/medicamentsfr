package medicamentsparser

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
)

var generiquesList []entities.GeneriqueList
var medsType map[int]string

func GeneriquesParser(medicaments *[]entities.Medicament, mMap *map[int]entities.Medicament) ([]entities.GeneriqueList, map[int]entities.Generique) {
	// allGeneriques: []Generique
	allGeneriques := makeGeneriques(nil)
	// Create a map of all the generiques to reduce algorithm complexity
	generiquesMap := make(map[int]entities.Generique)
	for i := range allGeneriques {
		generiquesMap[allGeneriques[i].Group] = allGeneriques[i]
	}

	// generiques file: [groupid]:[]cis of medicaments in the same group
	generiquesFile := generiqueFileToJSON()

	// The medsType is a map where the key are the medicament cis and the value is the
	// type of generique
	medsType = createMedicamentGeneriqueType()

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
			Medicaments: getMedicamentsInArray(v, medicaments, mMap),
		}

		generiques = append(generiques, current)
	}

	marshalledGeneriques, err := json.MarshalIndent(generiques, "", " ")
	if err != nil {
		log.Println("An error has occurred when marshalling generiques", err)
	}
	_ = os.WriteFile("src/GeneriquesFull.json", marshalledGeneriques, 0644)
	log.Println("GeneriquesFull.json created")
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

func getMedicamentsInArray(medicamentsIds []int, medicaments *[]entities.Medicament, medicamentMap *map[int]entities.Medicament) []entities.GeneriqueMedicament {
	var medicamentsArray []entities.GeneriqueMedicament

	for _, v := range medicamentsIds {
		if medicament, ok := (*medicamentMap)[v]; ok {
			generiqueComposition := createGeneriqueComposition(&medicament.Composition)
			generiqueMed := entities.GeneriqueMedicament{
				Cis:                 (medicament.Cis),
				Denomination:        (medicament.Denomination),
				FormePharmaceutique: (medicament.FormePharmaceutique),
				Type:                medsType[medicament.Cis],
				Composition:         generiqueComposition,
			}
			medicamentsArray = append(medicamentsArray, generiqueMed)
		}
	}

	return medicamentsArray
}
