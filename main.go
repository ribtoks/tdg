package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type arrayFlags []string

func (af *arrayFlags) String() string {
	return strings.Join(*af, " ")
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, value)
	return nil
}

var (
	includePatternsFlag arrayFlags
	srcRootFlag         = flag.String("root", "./", "Path to the the root of source code")
	helpFlag            = flag.Bool("help", false, "Show help")
	verboseFlag         = flag.Bool("verbose", false, "Be verbose")
)

func main() {
	err := parseFlags()
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	env := NewEnvironment(*srcRootFlag)
	td := NewToDoGenerator(*srcRootFlag, includePatternsFlag)
	err, comments := td.Generate()
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
	if !*verboseFlag {
		log.SetOutput(ioutil.Discard)
	}
	return nil
}
