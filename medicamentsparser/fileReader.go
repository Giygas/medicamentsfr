package medicamentsparser

import (
	"encoding/json"
	"log"
	"os"

	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
)


func specialitesFileToJSON() ([]entities.Specialite){
	fileData, err := os.ReadFile("src/Specialites.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	
	var specialites []entities.Specialite
	err = json.Unmarshal(fileData, &specialites)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	//Returns the full Specialites json as a slice of Specialites
	return specialites
}

func compositionFileToJSON() ([]entities.Composition){
	fileData, err := os.ReadFile("src/Compositions.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	
	var data []entities.Composition
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	//Returns the full Composition json as a slice of Composition
	return data
}

func conditionFileToJSON() ([]entities.Condition){
	fileData, err := os.ReadFile("src/Conditions.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	
	var data []entities.Condition
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	//Returns the full Condition json as a slice of Condition
	return data
}

func generiqueFileToJSON() ([]entities.Generique){
	fileData, err := os.ReadFile("src/Generiques.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	
	var data []entities.Generique
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	//Returns the full Generique json as a slice of Generique
	return data
}

func presentationFileToJSON() ([]entities.Presentation){
	fileData, err := os.ReadFile("src/Presentations.json")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	
	var data []entities.Presentation
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal("Error decoding JSON:", err)
	}

	//Returns the full Presentation json as a slice of Presentation
	return data
}