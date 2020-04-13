# gollup
![GitHub Actions](https://github.com/mpppk/gollup/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/mpppk/gollup/branch/master/graph/badge.svg)](https://codecov.io/gh/mpppk/gollup)
[![GoDoc](https://godoc.org/github.com/mpppk/gollup?status.svg)](https://godoc.org/github.com/mpppk/gollup)

gollup is a bundler for golang with tree-shaking.

## Caution
Most Go users do **not** need this tool.  
For example, this does not contribute to reducing the binary size.  
One of the few use cases is competitive programming, such as [atcoder](https://atcoder.jp), where only one file can be submitted.

## Installation

```shell script
$ go get github.com/mpppk/gollup
```

## Usage
### Simple example

```shell script
$ tree .
```

main.go:
```go
package main

import "fmt"

func main() {
	v := f()
	fmt.Println(v)
}
```

sub.go:
```go
package main

func f() int {
	return 42
}
```

```shell script
$ gollup | goimports > output.go
```

output.go:
```go
package main

import (
	"fmt"
)

func main() {
	v := f()
	fmt.Println(v)
}
func f() int {
	return 42
}
```

### multi package example

main.go
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

lib/lib.go:
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

output.go:
```go
package main

import (
        "fmt"
        "math"
)

const ANSWER = 42

func main() {
        fmt.Println(F1(), lib_F1())
}
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
```