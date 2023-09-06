package main

import (
	"log"
	"time"

	"github.com/giygas/medicamentsfr/medicamentsparser"
	"github.com/go-co-op/gocron"
)

func scheduleMedicaments() {
	s := gocron.NewScheduler(time.UTC)

	//s.Every(1).Day().At("06:00").Do(func() { medicaments = medicamentsparser.ParseAllMedicaments() })

	s.Every(1).Day().Do(func() {
		log.Println("Staring cron job")
		medicaments = medicamentsparser.ParseAllMedicaments()
		generiques = medicamentsparser.GeneriquesParser(&medicaments)
	})

}
