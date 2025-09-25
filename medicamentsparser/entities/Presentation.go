package entities

type Presentation struct {
	Cis                  int     `json:"cis"`
	Cip7                 int     `json:"cip7"`
	Libelle              string  `json:"libelle"`
	StatusAdministratif  string  `json:"statusAdministratif"`
	EtatComercialisation string  `json:"etatComercialisation"`
	DateDeclaration      string  `json:"dateDeclaration"`
	Cip13                int     `json:"cip13"`
	Agreement            string  `json:"agreement"`
	TauxRemboursement    string  `json:"tauxRemboursement"`
	Prix                 float32 `json:"prix"`
}
