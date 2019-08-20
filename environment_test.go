package main

import (
	"fmt"
	"testing"
)

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

var envTests = []struct {
	in  []string
	out []string
}{
	{[]string{"GIT_DIR=5"}, []string{}},
	{[]string{"GIT_DIR", "TEST=6"}, []string{"TEST=6"}},
	{[]string{"TEST=6", "GIT_DIR"}, []string{"TEST=6"}},
	{[]string{"TEST=6"}, []string{"TEST=6"}},
}

func TestSliceWithoutGitDir(t *testing.T) {
	for i, tt := range envTests {
		t.Run(fmt.Sprintf("test_%v", i), func(t *testing.T) {
			without := sliceWithoutGitDir(tt.in)
			if !Equal(without, tt.out) {
				t.Errorf("got %v, expected %v", without, tt.out)
			}
		})
	}
}
