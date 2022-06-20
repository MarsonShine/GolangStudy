package v2

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Undo []func()

func (undo *Undo) Add(f func()) {
	*undo = append(*undo, f)
}
func (undo *Undo) Undo() error {
	fs := *undo
	if len(fs) == 0 {
		return errors.New("No functions to undo")
	}
	index := len(fs) - 1
	if f := fs[index]; f != nil {
		f()
		fs[index] = nil
	}
	*undo = fs[:index]
	return nil
}

type IntSet struct {
	data map[int]bool
	undo Undo
}

func NewIntSet() IntSet {
	return IntSet{data: make(map[int]bool)}
}
func (set *IntSet) Add(x int) {
	if !set.Contains(x) {
		set.data[x] = true
		set.undo.Add(func() { set.Delete(x) })
	} else {
		set.undo.Add(nil)
	}
}
func (set *IntSet) Delete(x int) {
	if set.Contains(x) {
		delete(set.data, x)
		set.undo.Add(func() { set.Add(x) })
	} else {
		set.undo.Add(nil)
	}
}
func (set *IntSet) Undo() error {
	return set.undo.Undo()
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
