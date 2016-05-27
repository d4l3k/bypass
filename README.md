# bypass [![GoDoc](https://godoc.org/github.com/d4l3k/bypass?status.svg)](https://godoc.org/github.com/d4l3k/bypass) [![Build Status](https://travis-ci.org/d4l3k/bypass.svg?branch=master)](https://travis-ci.org/d4l3k/bypass) [![Coverage Status](https://coveralls.io/repos/github/d4l3k/bypass/badge.svg?branch=master)](https://coveralls.io/github/d4l3k/bypass?branch=master)

bypass is an incredibly unsafe Go package that allows you to inspect private
parts of go objects.

* Fetches elements of a channel without modifying it
* Allows reflect to be used on unexported fields

## Installation

```bash
$ go get -u github.com/d4l3k/bypass
```
Note: bypass only supports Go 1.5+.

## Example

```go
package main

import "github.com/d4l3k/bypass"

func main() {
  c := make(chan int, 10)
  c <- 1
  c <- 2
  c <- 3
  out := bypass.WrapChan(c).Elems().([]int)
  fmt.Printf("%#v\n", out)
  // Expected: []int{1, 2, 3}
}
```

## How to access private methods

https://sitano.github.io/2016/04/28/golang-private/

## License
Copyright (c) 2016 [Tristan Rice](https://fn.lc) <rice@fn.lc>

bypass is licensed under the MIT license. See the LICENSE file for more
information.

bypass.go and bypasssafe.go are borrowed from
[go-spew](https://github.com/davecgh/go-spew) and have a seperate copyright
notice.
