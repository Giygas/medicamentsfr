package entities

type Composition struct {
	Cis                   int    `json:"cis"`
	ElementParmaceutique  string `json:"elementPharmaceutique"`
	CodeSubstance         int    `json:"codeSubstance"`
	DenominationSubstance string `json:"denominationSubstance"`
	Dosage                string `json:"dosage"`
	ReferenceDosage       string `json:"referenceDosage"`
	NatureComposant       string `json:"natureComposant"`
}
