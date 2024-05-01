package validators

type Validator interface {
}

type validator struct {
}

func New() Validator {
	return &validator{}
}
