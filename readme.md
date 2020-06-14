# gollup
![GitHub Actions](https://github.com/mpppk/gollup/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/mpppk/gollup/branch/master/graph/badge.svg)](https://codecov.io/gh/mpppk/gollup)
[![GoDoc](https://godoc.org/github.com/mpppk/gollup?status.svg)](https://godoc.org/github.com/mpppk/gollup)

gollup is a bundler for golang with tree-shaking.

## Caution
Most Go users do **not** need this tool.  
For example, this does not contribute to reducing the binary size.  
One of the few use cases is competitive programming, such as [atcoder](https://atcoder.jp), where only one file can be submitted.

## Current status: under development
50+ codes I bundled using gollup are accepted on atcoder.  
However, this tool is not yet stable and may generate incorrect code.
I'm concerned that users will be penalized on the contest due to a bug in this tool.
Please try gollup on the past contest, check the behavior, and consider whether to use it on the contest.
When using it on a contest, I recommend that you also prepare other bundling methods such as [bundle](https://godoc.org/golang.org/x/tools/cmd/bundle) or manual.

Also, currently gollup does not support the following situations:
* [const with duplicate name between packages](https://github.com/mpppk/gollup/blob/07e7d57a766dd48efaf771ad115d4bfa3506a5e5/cmd/root_test.go#L152)
* [struct with duplicate name between packages](https://github.com/mpppk/gollup/blob/07e7d57a766dd48efaf771ad115d4bfa3506a5e5/cmd/root_test.go#L160)
* [nested/embedded struct](https://github.com/mpppk/gollup/blob/07e7d57a766dd48efaf771ad115d4bfa3506a5e5/cmd/root_test.go#L168)

## Installation

```shell script
$ go get github.com/mpppk/gollup
```

## Usage
### Simple example

```shell
$ tree .
.
├── main.go
└── sub.go
```

`main.go`:
```go
package main

import "fmt"

func main() {
	v := f()
	fmt.Println(v)
}
```

`sub.go`:
```go
package main

func f() int {
	return 42
}

// unusedFunc will not be included in bundled code because this is unused
func unusedFunc() {}
```

```shell script
$ gollup > output.go
```

`output.go`:
```go
package main

import (
	"fmt"
)

func f() int {
	return 42
}
func main() {
	v := f()
	fmt.Println(v)
}
```

### Multi package example

```shell
$ tree .
.
├── lib
│   └── lib.go
└── main.go
```

`main.go`:
```go
package main

import (
	"fmt"

	"github.com/mpppk/gollup/testdata/test2/lib"
)

const ANSWER = 42

func main() {
	fmt.Println(F1(), lib.F1())
}

func F1() int {
	return f()
}

func f() int {
	return ANSWER
}
```

`lib/lib.go`:
```go
package lib

import "math"

func F1() float64 {
	return f2()
}

func f2() float64 {
	return math.Sqrt(42)
}
```

```shell script
$ gollup ./lib . > output.go
```

`output.go`:
```go
package main

import (
        "fmt"
        "math"
)

const ANSWER = 42

func F1() int {
        return f()
}
func f() int {
        return ANSWER
}
func lib_F1() float64 {
        return lib_f2()
}
func lib_f2() float64 {
        return math.Sqrt(42)
}
func main() {
        fmt.Println(F1(), lib_F1())
}
```