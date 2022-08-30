package main

import (
	"fmt"
	"sort"
	"testing"

	"golang.org/x/exp/slices"
)

func TestBubbleSort(t *testing.T) {
	for _, length := range []int{1, 2, 4, 6, 17, 32, 800} {
		testname := fmt.Sprintf("sort-len-%d", length)
		t.Run(testname, func(t *testing.T) {
			// Test that our bubble sort works by comparing it to the built-in sort.
			ss := makeRandomStrings(length)
			ss2 := slices.Clone(ss)
			ss3 := slices.Clone(ss)
			ss4 := slices.Clone(ss)

			sort.Strings(ss)
			bubbleSortInterface(sort.StringSlice(ss2))
			bubbleSortGeneric(ss3)
			bubbleSortFunc(ss4, func(a, b string) bool { return a < b })

			for i := range ss {
				if ss[i] != ss2[i] {
					t.Fatalf("strings mismatch at %d; %s != %s", i, ss[i], ss2[i])
				}
				if ss[i] != ss3[i] {
					t.Fatalf("generic mismatch at %d; %s != %s", i, ss[i], ss3[i])
				}
				if ss[i] != ss4[i] {
					t.Fatalf("generic mismatch at %d; %s != %s", i, ss[i], ss4[i])
				}
			}
		})
	}
}

const N = 1_000

func BenchmarkSortStringInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStrings(N)
		b.StartTimer()
		bubbleSortInterface(sort.StringSlice(ss))
	}
}

func BenchmarkSortStringGeneric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStrings(N)
		b.StartTimer()
		bubbleSortGeneric(ss)
	}
}

func BenchmarkSortStringFunc(b *testing.B) {
	lessFunc := func(a, b string) bool { return a < b }
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ss := makeRandomStrings(N)
		b.StartTimer()
		bubbleSortFunc(ss, lessFunc)
	}
}
