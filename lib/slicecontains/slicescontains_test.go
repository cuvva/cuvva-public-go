package slicecontains_test

import (
	"testing"

	"github.com/cuvva/cuvva-public-go/lib/slicecontains"
)

func TestEmpty(t *testing.T) {
	var a []string
	var b []string

	if !slicecontains.SameStrings(a, b) {
		t.Error("slices should be equal")
	}

	a = []string{}
	b = []string{}

	if !slicecontains.SameStrings(a, b) {
		t.Error("slices should be equal")
	}
}

func TestSameStringsMatch(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"a", "c"}

	if slicecontains.SameStrings(a, b) {
		t.Error("slices should not be equal")
	}
}

func TestSameStringsSubset(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"a"}

	if slicecontains.SameStrings(a, b) {
		t.Error("slices should not be equal")
	}
}

func TestElem(t *testing.T) {
	a := []string{"a", "b"}
	b := "b"
	c := "c"

	if !slicecontains.Elem(a, b) {
		t.Error("b is contained within a")
	}

	if slicecontains.Elem(a, c) {
		t.Error("c is not contained within a")
	}
}

func TestElemsEmpty(t *testing.T) {
	var a []int
	var b []int

	if !slicecontains.SameElems(a, b) {
		t.Error("slices should be equal")
	}

	a = []int{}
	b = []int{}

	if !slicecontains.SameElems(a, b) {
		t.Error("slices should be equal")
	}
}

func TestSameElemsMatch(t *testing.T) {
	a := []int{0, 1}
	b := []int{0, 2}

	if slicecontains.SameElems(a, b) {
		t.Error("slices should not be equal")
	}
}

func TestSameElemsStringsSubset(t *testing.T) {
	a := []int{0, 1}
	b := []int{0}

	if slicecontains.SameElems(a, b) {
		t.Error("slices should not be equal")
	}
}
