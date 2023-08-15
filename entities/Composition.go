package entities

type Composition struct {
	Cis 										int			`json:"cis"`
	ElementParmaceutique 		string	`json:"elementPharmaceutique"`
	CodeSubstance 					int			`json:"codeSubstance"`
	VoiesAdministration 	string	`json:"VoiesAdministration"`
	Dosage 									string	`json:"dosage"`
	ReferenceDosage 				string	`json:"referenceDosage"`
	NatureComposant 				string	`json:"natureComposant"`
}