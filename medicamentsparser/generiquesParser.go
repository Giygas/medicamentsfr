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
	// generiques file: [groupid]:[]cis of medicaments in the same group
	//TODO: remove duplicates of groupid
	generiquesFile := generiqueFileToJSON()

	var generiques []entities.GeneriqueList
	for _, v := range allGeneriques {
		stringGroup := strconv.Itoa(v.Group)

		current := entities.GeneriqueList{
			GroupId:     v.Group,
			Libelle:     v.Libelle,
			Type:        v.Type,
			Medicaments: getMedicamentsInArray(generiquesFile[stringGroup], medicaments, medicamentMap),
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

func getMedicamentsInArray(medicamentsIds []int, medicaments *[]entities.Medicament, medicamentMap map[int]*entities.Medicament) []entities.Medicament {
	var medicamentsArray []entities.Medicament

	for _, v := range medicamentsIds {
		if medicament, ok := medicamentMap[v]; ok {
			medicamentsArray = append(medicamentsArray, *medicament)
		}
	}

	return medicamentsArray
}
