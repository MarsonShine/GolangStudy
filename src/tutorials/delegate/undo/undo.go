package undo

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type IntSet struct {
	data map[int]bool
}

func NewIntSet() IntSet {
	return IntSet{make(map[int]bool)}
}

func (set *IntSet) Add(x int) {
	set.data[x] = true
}
func (set *IntSet) Delete(x int) {
	delete(set.data, x)
}
func (set *IntSet) Contains(x int) bool {
	return set.data[x]
}
func (set *IntSet) String() string {
	if len(set.data) == 0 {
		return "{}"
	}
	ints := make([]int, 0, len(set.data))
	for i := range set.data {
		ints = append(ints, i)
	}
	sort.Ints(ints)
	parts := make([]string, 0, len(ints))
	for _, i := range ints {
		parts = append(parts, fmt.Sprint(i))
	}
	return "{" + strings.Join(parts, ",") + "}"
}

type UndoableIntSet struct {
	IntSet
	functions []func()
}

func NewUndoableIntSet() UndoableIntSet {
	return UndoableIntSet{NewIntSet(), nil}
}
func (set *UndoableIntSet) Add(x int) {
	if !set.Contains(x) {
		set.data[x] = true
		set.functions = append(set.functions, func() {
			set.Delete(x)
		})
	} else {
		set.functions = append(set.functions, nil)
	}
}
func (set *UndoableIntSet) Delete(x int) {
	if set.Contains(x) {
		set.data[x] = true
		set.functions = append(set.functions, func() {
			set.Add(x)
		})
	} else {
		set.functions = append(set.functions, nil)
	}
}
func (set *UndoableIntSet) Undo() error {
	if len(set.functions) == 0 {
		return errors.New("No functions to undo")
	}
	index := len(set.functions) - 1
	if f := set.functions[index]; f != nil {
		f()
		set.functions[index] = nil
	}
	set.functions = set.functions[:index]
	return nil
}
