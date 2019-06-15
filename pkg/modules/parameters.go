package modules

import (
	"encoding/json"
	"strconv"

	"github.com/pajbot/pajbot2/pkg"
)

func boolPtr(v bool) *bool {
	return &v
}

var nullBuffer = []byte("null")

var _ pkg.ModuleParameter = &floatParameter{}
var _ pkg.ModuleParameter = &boolParameter{}

type baseParameter struct {
	description string
}

func (b baseParameter) Description() string {
	return b.description
}

type floatParameter struct {
	baseParameter

	defaultValue *float32
	value        *float32
}

type parameterSpec struct {
	Description  string
	DefaultValue interface{}
}

func floatPtr(v float32) *float32 {
	return &v
}

func newFloatParameter(spec parameterSpec) *floatParameter {
	p := &floatParameter{}
	if spec.DefaultValue == nil {
		p.defaultValue = floatPtr(0.0)
	} else {
		var defaultValue float32
		var ok bool
		defaultValue, ok = spec.DefaultValue.(float32)
		if !ok {
			p.defaultValue = floatPtr(0.0)
		} else {
			p.defaultValue = &defaultValue
		}
	}

	p.description = spec.Description

	return p
}

func (p *floatParameter) DefaultValue() interface{} {
	return p.defaultValue
}

func (p *floatParameter) Get() float32 {
	if p.value != nil {
		return *p.value
	}

	if p.defaultValue != nil {
		return *p.defaultValue
	}

	return 0
}

func (p *floatParameter) Set(v float32) {
	p.value = &v
}

func (p *floatParameter) Reset() {
	p.value = nil
}

func (p *floatParameter) Parse(s string) error {
	if s == "reset" {
		p.Reset()
		return nil
	}

	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}

	p.Set(float32(v))

	return nil
}

func (p floatParameter) MarshalJSON() ([]byte, error) {
	if p.value != nil {
		return json.Marshal(p.value)
	}

	return nullBuffer, nil
}

func (p *floatParameter) UnmarshalJSON(b []byte) error {
	var v float32
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	p.value = &v

	return nil
}

type boolParameter struct {
	baseParameter

	defaultValue *bool
	value        *bool
}

func newBoolParameter(spec parameterSpec) *boolParameter {
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

func (p *boolParameter) Get() bool {
	if p.value != nil {
		return *p.value
	}

	if p.defaultValue != nil {
		return *p.defaultValue
	}

	return false
}

func (p *boolParameter) Set(v bool) {
	p.value = &v
}

func (p *boolParameter) Reset() {
	p.value = nil
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
