package modules

import (
	"encoding/json"
	"fmt"
	"log"
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

	hasBeenSet bool
}

func (b baseParameter) HasBeenSet() bool {
	return b.hasBeenSet
}

func (b baseParameter) Description() string {
	return b.description
}

type floatParameter struct {
	baseParameter

	defaultValue *float32

	value *float32
}

func (p *boolParameter) Get() interface{} {
	return p.value
}

func (p *floatParameter) Get() interface{} {
	return p.value
}

func (p *boolParameter) HasValue() bool {
	return p.value != nil
}

func (p *floatParameter) HasValue() bool {
	return p.value != nil
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
			log.Println("INVALID TYPE FOR THE DEFAULT VALUE", spec.DefaultValue)
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

func (p *floatParameter) Link(v interface{}) {
	linkedValue, ok := v.(*float32)
	if !ok {
		log.Println("Wrong value type!")
		return
	}

	p.value = linkedValue

	if p.defaultValue != nil {
		*p.value = *p.defaultValue
	}
}

func (p *floatParameter) Set(v float32) {
	if p.value == nil {
		log.Println("Set called on a float parameter without a value link")
		return
	}

	p.hasBeenSet = true

	*p.value = v
}

func (p *floatParameter) Reset() {
	var v float32
	if p.defaultValue != nil {
		v = *p.defaultValue
	}

	*p.value = v
	p.hasBeenSet = false
}

func (p *boolParameter) SetInterface(i interface{}) {
	// TODO: make some better type checks, maybe allow to set int?
	switch v := i.(type) {
	case bool:
		p.Set(v)
	case int:
		p.Set(v != 0)
	}
}

func (p *floatParameter) SetInterface(i interface{}) {
	// TODO: make some better type checks, maybe allow to set int?
	switch v := i.(type) {
	case float32:
		p.Set(v)
		return
	case float64:
		p.Set(float32(v))
		return
	case int:
		p.Set(float32(v))
		return
	}

	log.Printf("Unable to set value: %s (Type %T)", i, i)
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

func (p *floatParameter) MarshalJSON() ([]byte, error) {
	if p.value != nil {
		return json.Marshal(p.value)
	}

	return nullBuffer, nil
}

func (p *floatParameter) UnmarshalJSON(b []byte) error {
	fmt.Println("UNMARSHAL JSON:", string(b))
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
