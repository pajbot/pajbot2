package modules

import (
	"encoding/json"
	"strconv"
)

var nullBuffer = []byte("null")

type floatParameter struct {
	defaultValue *float32
	value        *float32
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
	defaultValue *bool
	value        *bool
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
