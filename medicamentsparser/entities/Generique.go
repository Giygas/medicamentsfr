package entities

type Generique struct {
	Cis     int    `json:"cis"`
	Group   int    `json:"group"`
	Libelle string `json:"libelle"`
	Type    string `json:"type"`
}
