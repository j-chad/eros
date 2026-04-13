package testutil

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

// Equal fails the test if got != want.
func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// NotEqual fails the test if got == want.
func NotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("got %v, want it to differ", got)
	}
}

// ElementsMatch fails the test if got and want do not contain the same elements
// (order-independent). Comparison uses the provided eq function.
func ElementsMatch[T comparable](t *testing.T, got, want []T) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("ElementsMatch: length mismatch: got %d, want %d", len(got), len(want))
		return
	}
	used := make([]bool, len(want))
	for i, g := range got {
		found := false
		for j, w := range want {
			if !used[j] && g == w {
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ElementsMatch: got[%d] = %v has no unused match in want", i, g)
			return
		}
	}
}

// NilErr fails the test if err is not nil.
func NilErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// NotNilErr fails the test if err is nil.
func NotNilErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

// ErrorIs fails the test if !errors.Is(err, target).
func ErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("got error %v, want errors.Is(%v)", err, target)
	}
}

// ErrorIsNot fails the test if errors.Is(err, target).
func ErrorIsNot(t *testing.T, err, target error) {
	t.Helper()
	if errors.Is(err, target) {
		t.Errorf("got error %v, want it to NOT match %v", err, target)
	}
}

// True fails the test if v is false.
func True(t *testing.T, v bool, msg ...string) {
	t.Helper()
	if !v {
		if len(msg) > 0 {
			t.Error(msg[0])
		} else {
			t.Error("expected true, got false")
		}
	}
}

// False fails the test if v is true.
func False(t *testing.T, v bool, msg ...string) {
	t.Helper()
	if v {
		if len(msg) > 0 {
			t.Error(msg[0])
		} else {
			t.Error("expected false, got true")
		}
	}
}

// Nil fails the test if v is not nil. Handles typed nils (e.g. (*T)(nil)).
func Nil(t *testing.T, v any) {
	t.Helper()
	if v == nil {
		return
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface ||
		rv.Kind() == reflect.Map || rv.Kind() == reflect.Slice ||
		rv.Kind() == reflect.Chan || rv.Kind() == reflect.Func {
		if rv.IsNil() {
			return
		}
	}
	t.Errorf("got %v (%T), want nil", v, v)
}

// NotNil fails the test if v is nil. Handles typed nils.
func NotNil(t *testing.T, v any) {
	t.Helper()
	if v == nil {
		t.Fatal("got nil, want non-nil")
		return
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface ||
		rv.Kind() == reflect.Map || rv.Kind() == reflect.Slice ||
		rv.Kind() == reflect.Chan || rv.Kind() == reflect.Func {
		if rv.IsNil() {
			t.Fatal("got typed nil, want non-nil")
		}
	}
}

// TypeAssert asserts that v is of type T and returns it. Fails the test otherwise.
func TypeAssert[T any](t *testing.T, v any) T {
	t.Helper()
	result, ok := v.(T)
	if !ok {
		t.Fatalf("type assertion failed: got %T, want %s", v, fmt.Sprintf("%T", *new(T)))
	}
	return result
}
