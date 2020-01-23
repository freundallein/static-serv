package server

import (
	"reflect"
	"testing"
)

func TestNewWithError(t *testing.T) {
	observed, err := New(nil)
	if err == nil {
		t.Error("expected", ErrNoOptions.Error(), "got", err.Error())
	}
	if observed != nil {
		t.Error("expected", nil, "got", observed)
	}
}
func TestNew(t *testing.T) {
	observed, err := New(&Options{})
	if err != nil {
		t.Error(err.Error())
	}
	observedType := reflect.TypeOf(observed)
	expectedType := reflect.TypeOf(&Server{})
	if observedType != expectedType {
		t.Error("Expected", expectedType, "got", observedType)
	}
}
