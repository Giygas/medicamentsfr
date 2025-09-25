// Package entities contains the entities used in the medicamentsparser package
package entities

type Condition struct {
	Cis       int    `json:"cis"`
	Condition string `json:"condition"`
}
