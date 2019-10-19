package modules

import (
	"encoding/json"
	"log"
)

type stringParameter struct {
	baseParameter

	defaultValue *string
	value        *string
}

func stringPtr(v string) *string {
	return &v
}

func newStringParameter(spec parameterSpec) *stringParameter {
	p := &stringParameter{}

	if spec.DefaultValue == nil {
		p.defaultValue = stringPtr("")
	} else {
		var defaultValue string
		var ok bool
		defaultValue, ok = spec.DefaultValue.(string)
		if !ok {
			p.defaultValue = stringPtr("")
		} else {
			p.defaultValue = &defaultValue
		}
	}

	p.description = spec.Description

	return p
}

func (p *stringParameter) Set(v string) {
	if p.value == nil {
		log.Println("Set called on a parameter without a value link")
		return
	}

	p.hasBeenSet = true

	*p.value = v
}

func (p *stringParameter) Reset() {
	var v string
	if p.defaultValue != nil {
		v = *p.defaultValue
	}

	*p.value = v
	p.hasBeenSet = false
}

func (p *stringParameter) Parse(s string) error {
	// XXX: This means the string parameter cannot be literally "reset". is this bad?
	if s == "reset" {
		p.Reset()
		return nil
	}

	p.Set(s)

	return nil
}

func (p stringParameter) MarshalJSON() ([]byte, error) {
	if p.value != nil {
		return json.Marshal(p.value)
	}

	return nullBuffer, nil
}

func (p *stringParameter) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	p.value = &v

	return nil
}

func (p *stringParameter) DefaultValue() interface{} {
	return p.defaultValue
}

func (p *stringParameter) SetInterface(i interface{}) {
	// TODO: make some better type checks, maybe allow to set int
	switch v := i.(type) {
	case string:
		p.Set(v)
	}
}

func (p *stringParameter) Link(v interface{}) {
	linkedValue, ok := v.(*string)
	if !ok {
		log.Println("Wrong value type!")
		return
	}

	p.value = linkedValue

	if p.defaultValue != nil {
		*p.value = *p.defaultValue
	}
}

func (p *stringParameter) HasValue() bool {
	return p.value != nil
}

func (p *stringParameter) Get() interface{} {
	return p.value
}
