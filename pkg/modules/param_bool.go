package modules

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/pajbot/pajbot2/pkg"
)

var _ pkg.ModuleParameter = &boolParameter{}

func boolPtr(v bool) *bool {
	return &v
}

type boolParameter struct {
	baseParameter

	defaultValue *bool
	value        *bool
}

func (p *boolParameter) Get() interface{} {
	return p.value
}

func (p *boolParameter) HasValue() bool {
	return p.value != nil
}

func (p *boolParameter) Link(v interface{}) {
	linkedValue, ok := v.(*bool)
	if !ok {
		log.Println("Wrong value type!")
		return
	}

	p.value = linkedValue

	if p.defaultValue != nil {
		*p.value = *p.defaultValue
	}
}

func (p *boolParameter) SetInterface(i interface{}) {
	// TODO: make some better type checks, maybe allow to set int
	switch v := i.(type) {
	case bool:
		p.Set(v)
	case int:
		p.Set(v != 0)
	}
}

func (p *boolParameter) Set(v bool) {
	if p.value == nil {
		log.Println("Set called on a parameter without a value link")
		return
	}

	p.hasBeenSet = true

	*p.value = v
}

func (p *boolParameter) Reset() {
	var v bool
	if p.defaultValue != nil {
		v = *p.defaultValue
	}

	*p.value = v
	p.hasBeenSet = false
}

func (p *boolParameter) Parse(s string) error {
	if s == "reset" {
		p.Reset()
		return nil
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	p.Set(v)

	return nil
}

func (p boolParameter) MarshalJSON() ([]byte, error) {
	if p.value != nil {
		return json.Marshal(p.value)
	}

	return nullBuffer, nil
}

func (p *boolParameter) UnmarshalJSON(b []byte) error {
	var v bool
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	p.value = &v

	return nil
}

func (p *boolParameter) DefaultValue() interface{} {
	return p.defaultValue
}

func NewBoolParameter(spec ParameterSpec) *boolParameter {
	p := &boolParameter{}

	if spec.DefaultValue == nil {
		p.defaultValue = boolPtr(false)
	} else {
		var defaultValue bool
		var ok bool
		defaultValue, ok = spec.DefaultValue.(bool)
		if !ok {
			p.defaultValue = boolPtr(false)
		} else {
			p.defaultValue = &defaultValue
		}
	}

	p.description = spec.Description

	return p
}
