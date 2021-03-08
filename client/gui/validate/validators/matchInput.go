package validators

import (
	"fmt"
	"gioui.org/widget"
	"github.com/skynet0590/atomicSwapTool/client/gui/validate"
	"strings"
)

type (
	matchInputValidator struct {
		matchedInput *widget.Editor
		fieldName    string
	}
)

func (v *matchInputValidator) Validate(input string) (valid bool, errTxt string) {
	if input != v.matchedInput.Text() {
		return false, v.errorTxt()
	}
	return true, ""
}

func (v *matchInputValidator) errorTxt() string {
	name := v.fieldName
	if name == "" {
		name = "field"
	}
	return fmt.Sprintf("Your %s is not matched", name)
}

func MatchedInput(fieldName string, input *widget.Editor) validate.Validator {
	return &matchInputValidator{
		fieldName:    strings.Trim(fieldName, " "),
		matchedInput: input,
	}
}
