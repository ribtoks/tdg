package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type arrayFlags []string

func (af *arrayFlags) String() string {
	return strings.Join(*af, " ")
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}

const (
	appName = "tdg"
)

var (
	includePatternsFlag arrayFlags
	srcRootFlag         = flag.String("root", "./", "Path to the the root of source code")
	helpFlag            = flag.Bool("help", false, "Show help")
	verboseFlag         = flag.Bool("verbose", false, "Output human-readable json")
	minWordCountFlag    = flag.Int("min-words", 3, "Skip comments with less than minimum words")
	minCharsFlag        = flag.Int("min-chars", 30, "Include comments with more chars than this")
	stdoutFlag          = flag.Bool("stdout", false, "Duplicate logs to stdout")
	logPathFlag         = flag.String("log", "tdg.log", "Path to the logfile")
)

func main() {
	err := parseFlags()
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	logfile, err := setupLogging()
	if err == nil {
		defer logfile.Close()
	}

	env := NewEnvironment(*srcRootFlag)
	td := NewToDoGenerator(*srcRootFlag, includePatternsFlag, *minWordCountFlag, *minCharsFlag)
	start := time.Now()
	comments, err := td.Generate()
	elapsed := time.Since(start)
	log.Printf("Generation took %s", elapsed)

	if err != nil {
		log.Fatal(err)
	}

	result := struct {
		Root     string         `json:"root"`
		Branch   string         `json:"branch"`
		Author   string         `json:"author"`
		Project  string         `json:"project"`
		Comments []*ToDoComment `json:"comments"`
	}{
		Root:     td.root,
		Branch:   env.Branch(),
		Author:   env.Author(),
		Project:  env.Project(),
		Comments: comments,
	}
	var js []byte
	if *verboseFlag {
		js, err = json.MarshalIndent(result, "", "  ")
	} else {
		js, err = json.Marshal(result)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(js))
}

func parseFlags() error {
	flag.Var(&includePatternsFlag, "include", "Include pattern (can be specified multiple times)")
	flag.Parse()
	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	srcRoot, err := os.Stat(*srcRootFlag)
	if os.IsNotExist(err) {
		return err
	}
	if !srcRoot.IsDir() {
		return errors.New("Root path does not point to a directory")
	}
	return nil
}

func setupLogging() (f *os.File, err error) {
	f, err = os.OpenFile(*logPathFlag, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		if *stdoutFlag {
			fmt.Printf("error opening file: %v", *logPathFlag)
		}
		return nil, err
	}

	if *stdoutFlag {
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	} else {
		log.SetOutput(f)
	}

	log.Println("------------------------------")
	log.Printf("%v log started", appName)

	return f, err
}
