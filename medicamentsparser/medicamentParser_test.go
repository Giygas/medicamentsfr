package medicamentsparser_test

import (
	"reflect"
	"testing"
	
	"github.com/giygas/medicamentsfr/medicamentsparser"
)

func TestParseAllMedicaments(t *testing.T) {
	result := medicamentsparser.ParseAllMedicaments()


	// Check if the result is of the expected type
	if reflect.TypeOf(result).String() != "[]entities.Medicament" {
		t.Errorf("Expected result to be of type []entities.Medicament, but got %v", reflect.TypeOf(result))
	}
	
}
