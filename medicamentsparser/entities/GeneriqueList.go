package entities

type GeneriqueList struct {
	GroupId     int                   `json:"groupId"`
	Libelle     string                `json:"libelle"`
	Medicaments []GeneriqueMedicament `json:"medicaments"`
}

type GeneriqueMedicament struct {
	Cis                 int    `json:"cis"`
	Denomination        string `json:"elementPharmaceutique"`
	FormePharmaceutique string `json:"formePharmaceutique"`
}
