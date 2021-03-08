package validators

import (
	"fmt"
	"github.com/skynet0590/atomicSwapTool/client/gui/validate"
	"strings"
)

type (
	requiredValidator struct {
		fieldName string
	}
)

func (v *requiredValidator) Validate(input string) (valid bool, errTxt string)  {
	if input == "" {
		return false, v.errorTxt()
	}
	return true, ""
}

func (v *requiredValidator) errorTxt() string {
	name := v.fieldName
	if name == "" {
		name = "The field"
	}
	return fmt.Sprintf("%s is required", name)
}

func Required(fieldName string) validate.Validator {
	return &requiredValidator{
		fieldName: strings.Trim(fieldName, " "),
	}
}
