package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type Environment struct {
	root       string
	branch     string
	author     string
	initBranch sync.Once
	initAuthor sync.Once
}

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
	}()
	return env
}

func (env *Environment) Run(cmd string, arg ...string) string {
	command := exec.Command(cmd, arg...)
	command.Dir = env.root
	out, err := command.Output()
	if err != nil {
		log.Print(err)
		return ""
	} else {
		return strings.TrimSpace(string(out))
	}
}

func (env *Environment) Branch() string {
	env.initBranch.Do(func() {
		env.branch = env.Run("git", "rev-parse", "--abbrev-ref", "HEAD")
	})
	return env.branch
}

func (env *Environment) Author() string {
	env.initAuthor.Do(func() {
		env.author = env.Run("git", "config", "user.name")
	})
	return env.author

}
