# `licenser` is a tool for maintenance of license header in source code

## Install

`go install go.bug.st/licenser`

## Usage

Create a `doc.go` with the header you want, for example:

```go
//
// Copyright 2022 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package main

// Other docs here...
```

run:

```bash
$ licenser .
Golang project detected
Extracting license from doc.go
IGNORED: LICENSE
IGNORED: README.md
OK doc.go
IGNORED: go.mod
IGNORED: go.sum
OK main.go
```

the tool will automatically copy the license found in `doc.go` to all your other golang source files. If the license is changed, the license on the other files will be updated as well.

WARNING: Always use this tool in a source tree that has a version control system like `git` in place!

This tool will update files **in-place**, this means that the old content will be **overwritten without confirmation**, so better to have a backup in case of failures.
