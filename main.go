package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var files = map[string]string {
	"specialites": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_bdpm.txt",
	"presentations": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_CIP_bdpm.txt",
	"compositions": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_COMPO_bdpm.txt",
	"generiques": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_GENER_bdpm.txt",
	"conditions": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_CPD_bdpm.txt",
}

func main() {
	
	//Create the files directory if it doesn't exists
	filePath := filepath.Join(".", "files")
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic( err)
	}
	
	//Download all the files and convert from windows 1252 to utf8
	downloadAndParseAll(files)
	
	//Pass only the names of the files to the function
	var filesNames []string
	for name := range(files) {
		filesNames = append(filesNames, name)
	}
	fmt.Println(filesNames)
	
	// TODO:Parse to JSON each downloaded file
	makePresentations()
	
}
