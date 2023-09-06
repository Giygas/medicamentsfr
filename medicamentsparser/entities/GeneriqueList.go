package entities

type GeneriqueList struct {
	GroupId     int          `json:"groupid"`
	Libelle     string       `json:"libelle"`
	Type        string       `json:"type"`
	Medicaments []Medicament `json:"medicaments"`
}
