package entities

type GeneriqueList struct {
	GroupId     int                   `json:"groupid"`
	Libelle     string                `json:"libelle"`
	Type        string                `json:"type"`
	Medicaments []GeneriqueMedicament `json:"medicaments"`
}

type GeneriqueMedicament struct {
	Cis                 int    `json:"cis"`
	Denomination        string `json:"elementPharmaceutique"`
	FormePharmaceutique string `json:"formePharmaceutique"`
}
