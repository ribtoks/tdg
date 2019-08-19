package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/zieckey/goini"
)

var (
	commentPrefixes = [...]string{"TODO: ", "FIXME: ", "BUG: ", "HACK: "}
	emptyRunes      = [...]rune{}
	categoryIniKey  = "category"
	issueIniKey     = "issue"
)

// ToDoComment a task that is parsed from TODO comment
type ToDoComment struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Issue    int    `json:"issue,omitempty"`
	Category string `json:"category,omitempty"`
}

// ToDoGenerator is responsible for parsing code base to ToDoComments
type ToDoGenerator struct {
	root         string
	filters      []*regexp.Regexp
	commentsChan chan *ToDoComment
	commentsWG   sync.WaitGroup
	comments     []*ToDoComment
	minWords     int
}

// NewToDoGenerator creates new generator for a source root
func NewToDoGenerator(root string, filters []string, minWords int) *ToDoGenerator {
	log.Printf("Using %v filters", filters)
	rfilters := make([]*regexp.Regexp, 0, len(filters))
	for _, f := range filters {
		rfilters = append(rfilters, regexp.MustCompile(f))
	}
	absolutePath, err := filepath.Abs(root)
	if err != nil {
		absolutePath = root
	}
	td := &ToDoGenerator{
		root:         absolutePath,
		filters:      rfilters,
		minWords:     minWords,
		commentsChan: make(chan *ToDoComment),
		comments:     make([]*ToDoComment, 0),
	}
	go td.processComments()
	return td
}

// Generate is an entry point to comment generation
func (td *ToDoGenerator) Generate() ([]*ToDoComment, error) {
	matchesCount := 0
	err := filepath.Walk(td.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		anyMatch := false
		for _, f := range td.filters {
			if f.MatchString(path) {
				anyMatch = true
				break
			}
		}
		if !anyMatch && len(td.filters) > 0 {
			return nil
		}

		matchesCount++
		td.commentsWG.Add(1)
		go td.parseFile(path)

		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("Matched files: %v", matchesCount)
	td.commentsWG.Wait()
	close(td.commentsChan)
	return td.comments, nil
}

func countTitleWords(s string) int {
	words := strings.Fields(s)
	count := 0
	for _, w := range words {
		if len(w) > 2 {
			count++
		}
	}
	return count
}

func (td *ToDoGenerator) processComments() {
	for c := range td.commentsChan {
		if countTitleWords(c.Title) >= td.minWords {
			td.comments = append(td.comments, c)
		}
		td.commentsWG.Done()
	}
}

func isCommentRune(r rune) bool {
	return r == '/' ||
		r == '#' ||
		r == '%' ||
		r == ';' ||
		r == '*'
}

// try to parse comment body from commented line
func parseComment(line string) []rune {
	runes := []rune(line)
	i := 0
	size := len(runes)
	// skip prefix whitespace
	for i < size && unicode.IsSpace(runes[i]) {
		i++
	}
	hasComment := false
	// skip comment symbols themselves
	for i < size && isCommentRune(runes[i]) {
		i++
		hasComment = true
	}
	if !hasComment {
		return nil
	}
	// and skip space again
	for i < size && unicode.IsSpace(runes[i]) {
		i++
	}
	j := size - 1
	// skip suffix whitespace
	for j > i && unicode.IsSpace(runes[j]) {
		j--
	}
	// empty comment
	if i >= size || j < 0 || i >= j {
		return emptyRunes[:]
	}
	return runes[i : j+1]
}

func startsWith(s, pr []rune) bool {
	// do not check length (it's checked above)
	for i, p := range pr {
		if unicode.ToUpper(s[i]) != p {
			return false
		}
	}
	return true
}

func parseToDoTitle(line []rune) (ctype, title []rune) {
	if line == nil || len(line) == 0 {
		return nil, nil
	}
	size := len(line)
	for _, pr := range commentPrefixes {
		prlen := len(pr)
		if size > prlen && startsWith(line, []rune(pr)) {
			// without last ':<space>'
			ctype = []rune(pr)[:prlen-2]
			title = line[prlen:]
			return
		}
	}

	return nil, nil
}

// NewComment creates new task from parsed comment lines
func NewComment(path string, lineNumber int, ctype string, body []string) *ToDoComment {
	if body == nil || len(body) == 0 {
		return nil
	}
	var commentBody string
	var issue int
	var category string

	if len(body) > 1 {
		if strings.Contains(body[1], "=") {
			ini := goini.New()
			err := ini.Parse([]byte(body[1]), " ", "=")
			if err == nil {
				if v, ok := ini.Get(categoryIniKey); ok {
					category = v
				}
				if v, ok := ini.Get(issueIniKey); ok {
					if i, err := strconv.Atoi(v); err == nil {
						issue = i
					}
				}
			} else {
				log.Print(err)
			}
		}
		if len(category) > 0 || issue > 0 {
			commentBody = strings.Join(body[2:], "\n")
		} else {
			commentBody = strings.Join(body[1:], "\n")
		}
		commentBody = strings.TrimSpace(commentBody)
	}
	return &ToDoComment{
		Type:     string(ctype),
		Title:    body[0],
		Body:     commentBody,
		File:     path,
		Line:     lineNumber,
		Category: category,
		Issue:    issue,
	}
}

func (td *ToDoGenerator) accountComment(path string, lineNumber int, ctype string, body []string) {

	relativePath, err := filepath.Rel(td.root, path)
	if err != nil {
		relativePath = path
	}
	c := NewComment(relativePath, lineNumber, ctype, body)
	td.commentsWG.Add(1)
	go func(cmt *ToDoComment) {
		td.commentsChan <- cmt
	}(c)
}

func (td *ToDoGenerator) parseFile(path string) {
	defer td.commentsWG.Done()
	f, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var todo []string
	var lastType string
	var lastStart int
	lineNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		if c := parseComment(line); c != nil {
			// current comment is new TODO-like commment
			if ctype, title := parseToDoTitle(c); title != nil {
				// do we need to finalize previous
				if lastType != "" {
					td.accountComment(path, lastStart, lastType, todo)
				}
				// construct new one
				lastType = string(ctype)
				lastStart = lineNumber - 1
				todo = make([]string, 0)
				todo = append(todo, string(title))
			} else if lastType != "" {
				// continue consecutive comment line
				todo = append(todo, string(c))
			}
		} else {
			// not a comment anymore: finalize
			if lastType != "" {
				td.accountComment(path, lastStart, lastType, todo)
				lastType = ""
			}
		}
	}
	// detect todo item at the end of the file
	if lastType != "" {
		td.accountComment(path, lastStart, lastType, todo)
	}
}
