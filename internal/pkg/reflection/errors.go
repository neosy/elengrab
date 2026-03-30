package reflection

import "errors"

var (
	ErrFirstArgumentTypeMustPointerStructure       = errors.New("the first argument must be a pointer to a structure")
	ErrSecondArgumentTypeMustPointerFieldStructure = errors.New("the second argument must be a pointer to a field of the structure")
	ErrStructureFieldNameEmpty                     = errors.New("empty structure field name")
	ErrFieldNotFound                               = errors.New("field not found")
	ErrInputMustBePointerToStruct                  = errors.New("input must be a pointer to a struct")
)
