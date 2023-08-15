package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)


var files = map[string]string {
	"Specialites": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_bdpm.txt",
	"Presentations": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_CIP_bdpm.txt",
	"Compositions": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_COMPO_bdpm.txt",
	"Generiques": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_GENER_bdpm.txt",
	"Conditions": "https://base-donnees-publique.medicaments.gouv.fr/telechargement.php?fichier=CIS_CPD_bdpm.txt",
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc / 1024 / 1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc / 1024 / 1024)
	fmt.Printf("\tSys = %v MiB", m.Sys / 1024 / 1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}


func main() {
	start := time.Now()
	//Create the files directory if it doesn't exists
	filePath := filepath.Join(".", "files")
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic( err)
	}
	filePath = filepath.Join(".", "src")
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		panic( err)
	}
	
	//Download all the files and convert from windows 1252 to utf8
	downloadAndParseAll(files)
	
	//Make all the json files concurrently
	var wg sync.WaitGroup
	wg.Add(5)
	
	go makePresentations(&wg)
	go makeSpecialites(&wg)
	go makeGeneriques(&wg)
	go makeCompositions(&wg)
	go makeConditions(&wg)
	
	wg.Wait()
	
	medicaments := parseAllMedicaments()

	fmt.Println(medicaments[0])
	
	timeElapsed := time.Since(start)
	fmt.Printf("The full database upgrade took: %s", timeElapsed)
	fmt.Println()
	PrintMemUsage()
}
