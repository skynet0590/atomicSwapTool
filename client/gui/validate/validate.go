package validate

type Validator interface {
	Validate(input string) (valid bool, errTxt string)
}
