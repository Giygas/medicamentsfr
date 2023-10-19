package main

import (
	"fmt"
	"log"
	"time"

	"github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/giygas/medicamentsfr/medicamentsparser/entities"
	"github.com/go-co-op/gocron"
)

func checkMedicaments(medicaments *[]entities.Medicament) {
	if len(*medicaments) == 0 {
		fmt.Println("medicaments slice is empty")
		return
	}

	for i, medicament := range *medicaments {
		if medicament.Cis == 0 {
			fmt.Printf("medicament at index %d has Cis set to 0\n", i)
		}
	}
}

func scheduleMedicaments(medicaments *[]entities.Medicament, mMap *map[int]entities.Medicament, generiques *[]entities.GeneriqueList, gMap *map[int]entities.Generique) {
	s := gocron.NewScheduler(time.Local)

<<<<<<< HEAD:apimedicaments/scheduler.go
	_, err := s.Every(1).Days().Do(func() {
=======
	_, err := s.Every(1).Days().At("06:00;18:00").StartImmediately().Do(func() {
>>>>>>> droplet-deploy:scheduler.go
		fmt.Println("Starting update at: ", time.Now())
		// Get the current time for calculating the total database update time
		start := time.Now()

		*medicaments = medicamentsparser.ParseAllMedicaments()
		checkMedicaments(medicaments)

		// Create a map of all medicaments to reduce algorithm complexity
		for i := range *medicaments {
			medicamentsMap[(*medicaments)[i].Cis] = (*medicaments)[i]
		}

		*generiques, generiquesMap = medicamentsparser.GeneriquesParser(medicaments, &medicamentsMap)
		fmt.Println("Finished update at: ", time.Now())

		// Total time updating
		timeElapsed := time.Since(start)
		fmt.Printf("The full database upgrade took: %s\n", timeElapsed)
	})
	if err != nil {
		log.Fatal("An error has ocurred while starting the chron job", err)
	}
	s.StartAsync()
}
