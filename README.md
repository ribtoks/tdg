# tdg
TODO get - tool to extract todo tasks from source code

[![Build Status](https://travis-ci.org/ribtoks/tdg.svg?branch=master)](https://travis-ci.org/ribtoks/tdg)
[![Build status](https://ci.appveyor.com/api/projects/status/ohve71khsdtv3ey6?svg=true)](https://ci.appveyor.com/project/Ribtoks/tdg)
[![Go Report Card](https://goreportcard.com/badge/github.com/ribtoks/tdg)](https://goreportcard.com/report/github.com/ribtoks/tdg)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/712b5193d6564beb88ba1e66ac1e0792)](https://www.codacy.com/app/ribtoks/tdg)
[![Maintainability](https://api.codeclimate.com/v1/badges/89dad5db195c7b5d3e90/maintainability)](https://codeclimate.com/github/ribtoks/tdg/maintainability)

![license](https://img.shields.io/badge/license-Apache_2.0-blue.svg)
![copyright](https://img.shields.io/badge/%C2%A9-Taras_Kushnir-blue.svg)
![language](https://img.shields.io/badge/language-go-blue.svg)

## About

This tool generates json from comments contained in the source code. Main use-case for it is to create automatic issues based on TODO/FIXME/BUG/HACK comments. This tool supports additional tag information in the comment (Category, Issue etc.).

Example of the comment:

    // TODO: This is title of the issue to create
    // category=SomeCategory issue=123
    // This is a multiline description of the issue
    // that will be in the "Body" property of the comment

Sample generated json (if run in this repository root as `tdg -verbose`):

    {
      "root": "C:\\Users\\latis\\go\\src\\github.com\\ribtoks\\tdg",
      "branch": "master",
      "author": "Taras Kushnir",
      "project": "tdg",
      "comments": [
        {
          "type": "TODO",
          "title": "This is title of the issue to create",
          "body": "This is a multiline description of the issue\nthat will be in the \"Body\" property of the comment",
          "file": "README.md",
          "line": 19,
          "issue": 123,
          "category": "SomeCategory"
        }
      ]
    }

Supported comments: `//`, `/*`, `#`, `%`, `;;` (adding new supported comments is trivial).

## Install

As simple as

    go get github.com/ribtoks/tdg

## Build

As simple as

    go build

## Usage

    -help
        Show help
    -include value
        Include pattern (can be specified multiple times)
    -log string
        Path to the logfile (default "tdg.log")
    -min-words int
        Skip comments with less than minimum words (default 3)
    -root string
        Path to the the root of source code (default "./")
    -stdout
        Duplicate logs to stdout
    -verbose
        Output human-readable json

Example:

    tdg -root ~/Projects/xpiks-root/xpiks/src/ -include "\.(cpp|h)$" -verbose

Include pattern is a regexp. With verbose flag you get human-readable json and log output in stdout. Without verbose flag this tool could be used as input for smth else like `curl`.

## How to contribute

-   [Fork](http://help.github.com/forking/) linuxdeploy repository on GitHub
-   Clone your fork locally
-   Configure the upstream repo (`git remote add upstream git@github.com:ribtoks/tdg.git`)
-   Create local branch (`git checkout -b your_feature`)
-   Work on your feature
-   Build and Run tests (`go tests -v`)
-   Push the branch to GitHub (`git push origin your_feature`)
-   Send a [pull request](https://help.github.com/articles/using-pull-requests) on GitHub
