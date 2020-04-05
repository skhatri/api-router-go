package test

import (
	"testing"
)

func NotNull(t *testing.T, in interface{}) {
	if in == nil {
		t.Fail()
	}
}

func Null(t *testing.T, in interface{}) {
	if in != nil {
		t.Fail()
	}
}

func EqualTo(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Logf("expected %v, actual %v", expected, actual)
		t.Fail()
	}
}
