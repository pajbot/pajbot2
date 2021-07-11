package modules

var nullBuffer = []byte("null")

type baseParameter struct {
	description string

	hasBeenSet bool
}

func (b baseParameter) HasBeenSet() bool {
	return b.hasBeenSet
}

func (b baseParameter) Description() string {
	return b.description
}

type ParameterSpec struct {
	Description  string
	DefaultValue interface{}
}
