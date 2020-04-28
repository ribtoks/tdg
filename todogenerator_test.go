package main

import (
	"fmt"
	"math"
	"testing"
)

const float64EqualityThreshold = 1e-9

var commentTests = []struct {
	in  string
	out string
}{
	{"//", ""},
	{"#", ""},
	{"%  ", ""},
	{"// test contents\t", "test contents"},
	{"# test contents", "test contents"},
	{"% test contents", "test contents"},
	{"/* test contents ", "test contents"},
	{"/** test contents ", "test contents"},
	{"// TODO: test contents", "TODO: test contents"},
	{"//TODO: test contents\t", "TODO: test contents"},
	{"//TODO(author): test contents", "TODO(author): test contents"},
}

func TestParseComment(t *testing.T) {
	for _, tt := range commentTests {
		t.Run(tt.in, func(t *testing.T) {
			comment := parseComment(tt.in)
			if comment == nil {
				t.Errorf("got nil, expected %v", tt.out)
			}
			if string(comment) != tt.out {
				t.Errorf("got %v, expected %v", comment, tt.out)
			}
		})
	}
}

var startsWithTests = []struct {
	in     string
	prefix string
	out    bool
}{
	{"TODO: test code", "TODO", true},
	{"TODO(author): test code", "TODO", true},
	{"TODO: test code", "HACK", false},
	{"ToDo: test code", "TODO", true},
	{"todo: test code", "TODO", true},
	{"BUG: test code", "BUG", true},
	{"BUG:test code", "BUG", true},
	{"BUG(author):test code", "BUG", true},
}

func TestStartsWith(t *testing.T) {
	for _, tt := range startsWithTests {
		t.Run(tt.in, func(t *testing.T) {
			if startsWith([]rune(tt.in), []rune(tt.prefix)) != tt.out {
				t.Errorf("Test(%v): got %v, expected %v", tt.in, !tt.out, tt.out)
			}
		})
	}
}

var todoLineTests = []struct {
	in    string
	ctype string
	title string
}{
	{"TODO: test code", "TODO", "test code"},
	{"TODO test code", "TODO", "test code"},
	{"TODOtest code", "", ""},
	{"FIXME: test code", "FIXME", "test code"},
	{"BUG: test code", "BUG", "test code"},
	{"BUG  test code", "BUG", "test code"},
	{"HACK: test code", "HACK", "test code"},
	{"HACK:test code", "HACK", "test code"},
	{"TODO(author): test code", "TODO", "test code"},
	{"HACK(author): test code", "HACK", "test code"},
	{"BUG(author): test code", "BUG", "test code"},
}

func TestTodoTitle(t *testing.T) {
	for _, tt := range todoLineTests {
		t.Run(tt.in, func(t *testing.T) {
			ctype, title := parseToDoTitle([]rune(tt.in))
			if string(ctype) != tt.ctype {
				t.Errorf("Test(%v): got %v, expected %v", tt.in, string(ctype), tt.ctype)
			}
			if string(title) != tt.title {
				t.Errorf("Test(%v): got %v, expected %v", tt.in, string(title), tt.title)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

var commentConstructorTests = []struct {
	in       []string
	title    string
	body     string
	category string
	issue    int
	estimate float64
}{
	{[]string{"issue title"}, "issue title", "", "", 0, 0.0},
	{[]string{"issue title", "category=Test"}, "issue title", "", "Test", 0, 0.0},
	{[]string{"issue title", "issue=123"}, "issue title", "", "", 123, 0.0},
	{[]string{"issue title", "estimate=30m"}, "issue title", "", "", 0, 0.5},
	{[]string{"issue title", "estimate=30x"}, "issue title", "estimate=30x", "", 0, 0.0},
	{[]string{"issue title", "estimate=30h"}, "issue title", "", "", 0, 30},
	{[]string{"issue title", "  category=Test issue=123 estimate=60m "}, "issue title", "", "Test", 123, 1.0},
	{[]string{"issue title", "issue=123", "third line"}, "issue title", "third line", "", 123, 0.0},
	{[]string{"issue title", "second line"}, "issue title", "second line", "", 0, 0.0},
	{[]string{"issue title", "second line", "third line"}, "issue title", "second line\nthird line", "", 0, 0.0},
}

func TestNewComment(t *testing.T) {
	for i, tt := range commentConstructorTests {
		t.Run(fmt.Sprintf("test_%v", i), func(t *testing.T) {
			c := NewComment("/path/", 0, "ctype", tt.in)
			if c.Title != tt.title {
				t.Errorf("Title error: got %v, expected %v", c.Title, tt.title)
			}
			if c.Body != tt.body {
				t.Errorf("Body error: got %v, expected %v", c.Body, tt.body)
			}
			if c.Category != tt.category {
				t.Errorf("Category error: got %v, expected %v", c.Category, tt.category)
			}
			if c.Issue != tt.issue {
				t.Errorf("Issue error: got %v, expected %v", c.Issue, tt.issue)
			}
			if !almostEqual(c.Estimate, tt.estimate) {
				t.Errorf("Estimate error: got %v, expected %v", c.Estimate, tt.estimate)
			}
		})
	}
}

var titleTests = []struct {
	in  string
	out int
}{
	{"a a a", 0},
	{"aaa", 1},
	{"aaa bbb c d", 2},
	{"aa bb cd", 0},
}

func TestTitleWords(t *testing.T) {
	for _, tt := range titleTests {
		t.Run(tt.in, func(t *testing.T) {
			wc := countTitleWords(tt.in)
			if wc != tt.out {
				t.Errorf("got %v, expected %v", wc, tt.out)
			}
		})
	}
}

var estimateTests = []struct {
	in  string
	out float64
}{
	{"", 0},
	{"30", 30},
	{"30h", 30},
	{"30m", 0.5},
	{"30x", 0},
}

func TestEstimates(t *testing.T) {
	for _, tt := range estimateTests {
		t.Run(tt.in, func(t *testing.T) {
			e, _ := parseEstimate(tt.in)
			if e != tt.out {
				t.Errorf("got %v, expected %v", e, tt.out)
			}
		})
	}
}
