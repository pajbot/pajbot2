package eventemitter

import (
	"errors"
	"testing"
)

var testError = errors.New("test")

func TestBadListen(t *testing.T) {
	e := New()
	_, err := e.Listen("asd", nil, 100)
	assertErrorsEqual(t, ErrBadCallback, err)
}

func TestGoodArgumentListen(t *testing.T) {
	e := New()
	_, err := e.Listen("asd", func(arguments map[string]interface{}) error {
		return nil
	}, 100)
	assertErrorsEqual(t, nil, err)
}

func TestGoodArgumentlessListen(t *testing.T) {
	e := New()
	_, err := e.Listen("asd", func() error {
		return nil
	}, 100)
	assertErrorsEqual(t, nil, err)
}

func TestEmitEmptyListener(t *testing.T) {
	e := New()
	n, err := e.Emit("asd", nil)
	assertIntsEqual(t, 0, n)
	assertErrorsEqual(t, nil, err)
}

func TestEmitOneListener(t *testing.T) {
	e := New()
	_, err := e.Listen("asd", func(arguments map[string]interface{}) error {
		return nil
	}, 100)
	assertErrorsEqual(t, nil, err)
	n, err := e.Emit("asd", nil)
	assertIntsEqual(t, 1, n)
	assertErrorsEqual(t, nil, err)
}

func TestEmitOneArgumentlessListener(t *testing.T) {
	e := New()
	_, err := e.Listen("asd", func() error {
		return nil
	}, 100)
	assertErrorsEqual(t, nil, err)
	n, err := e.Emit("asd", nil)
	assertIntsEqual(t, 1, n)
	assertErrorsEqual(t, nil, err)
}

func TestDisconnect(t *testing.T) {
	e := New()
	conn, err := e.Listen("asd", func() error {
		return nil
	}, 100)
	conn.Disconnected = true
	assertErrorsEqual(t, nil, err)
	n, err := e.Emit("asd", nil)
	assertIntsEqual(t, 0, n)
	assertErrorsEqual(t, nil, err)
}

func TestMultipleListeners(t *testing.T) {
	e := New()
	_, err := e.Listen("asd", func() error {
		return nil
	}, 100)
	assertErrorsEqual(t, nil, err)
	_, err = e.Listen("asd", func() error {
		return nil
	}, 100)
	assertErrorsEqual(t, nil, err)
	n, err := e.Emit("asd", nil)
	assertIntsEqual(t, 2, n)
}

func TestListenerError(t *testing.T) {
	e := New()
	_, err := e.Listen("asd", func() error {
		return testError
	}, 100)
	assertErrorsEqual(t, nil, err)
	_, err = e.Listen("asd", func() error {
		return nil
	}, 100)
	assertErrorsEqual(t, nil, err)
	n, err := e.Emit("asd", nil)
	assertIntsEqual(t, 0, n)
}

func TestListenerErrorArguments(t *testing.T) {
	e := New()
	_, err := e.Listen("asd", func(arguments map[string]interface{}) error {
		return testError
	}, 100)
	assertErrorsEqual(t, nil, err)
	_, err = e.Listen("asd", func(arguments map[string]interface{}) error {
		return nil
	}, 100)
	assertErrorsEqual(t, nil, err)
	n, err := e.Emit("asd", nil)
	assertIntsEqual(t, 0, n)
	assertErrorsEqual(t, testError, err)
}
