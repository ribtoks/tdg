package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Environment contains information about git repository
type Environment struct {
	root        string
	branch      string
	author      string
	project     string
	initBranch  sync.Once
	initAuthor  sync.Once
	initProject sync.Once
}

// NewEnvironment creates new instance of Environment struct
func NewEnvironment(root string) *Environment {
	absolutePath, err := filepath.Abs(root)
	if err != nil {
		absolutePath = root
	}
	env := &Environment{
		root: absolutePath,
	}
	go func() {
		log.Printf("Current branch is %v", env.Branch())
		log.Printf("Current author is %v", env.Author())
		log.Printf("Current project is %v", env.Project())
	}()
	return env
}

// Run executes a command in the environment's root
func (env *Environment) Run(cmd string, arg ...string) string {
	command := exec.Command(cmd, arg...)
	command.Dir = env.root
	out, err := command.Output()
	if err != nil {
		log.Print(err)
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Branch returns current git branch
func (env *Environment) Branch() string {
	env.initBranch.Do(func() {
		env.branch = env.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	})
	return env.branch
}

// Author returns current git author
func (env *Environment) Author() string {
	env.initAuthor.Do(func() {
		env.author = env.Run("git", "config", "user.name")
	})
	return env.author
}

// Project returns current git project name
func (env *Environment) Project() string {
	env.initProject.Do(func() {
		project := env.Run("git", "rev-parse", "--show-toplevel")
		env.project = env.Run("basename", project)
	})
	return env.project
}
