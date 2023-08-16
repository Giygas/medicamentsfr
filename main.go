package main

import (
	"fmt"
	"medicamentsfr/entities"
	"runtime"
	"time"
)

var medicaments []entities.Medicament


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

	//Download all the files and convert from windows 1252 to utf8
	downloadAndParseAll()
	
	//Parse all the downloaded files and create the medicaments.json
	medicaments = parseAllMedicaments()

	timeElapsed := time.Since(start)
	fmt.Printf("The full database upgrade took: %s", timeElapsed)
	
}
