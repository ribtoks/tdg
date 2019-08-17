# tdg
TODO get - tool to extract todo tasks from source code

[![Build Status](https://travis-ci.org/ribtoks/tdg.svg?branch=master)](https://travis-ci.org/ribtoks/tdg)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/712b5193d6564beb88ba1e66ac1e0792)](https://www.codacy.com/app/ribtoks/tdg)

## About

This tool generates json from comments contained in the source code. Main use-case for it is to create automatic issues based on TODO/FIXME/BUG/HACK comments. This tool supports additional tag information in the comment

Example of the comment:

    // TODO: This is title of the issue to create
    // [optional] category=SomeCategory issue=123
    // This is a multiline description of the issue
    // that will be in the "Body" property of the comment

Sample generated json:

    {
        "Root": "/Users/user/Projects/tdg",
        "Branch": "master",
        "Author": "Taras Kushnir",
        "Comments": [
            {
                "Type": "TODO",
                "Title": "This is title of the issue to create"
                "Body": "This is a multiline description of the issue\nthat will be in the \"Body\" property of the comment",
                "File": "main.go",
                "Line": 92,
                "Issue": 123,
                "Category": "SomeCategory"
            }
        ]
    }

## Build

As simple as

    go build

## Usage

    -help
            Show help
    -include value
            Include pattern (can be specified multiple times)
    -root string
            Path to the the root of source code (default "./")
    -verbose
            Be verbose

Example:

    ./tdg --root ~/Projects/xpiks-root/xpiks/src/ --include "\.cpp$" --include "\.h$" -verbose

Include pattern is a regexp. With verbose flag you get human-readable json and log output in stdout. Without verbose flag this tool could be used as input for smth else like `curl`.
