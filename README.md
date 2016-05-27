# bypass [![GoDoc](https://godoc.org/github.com/d4l3k/bypass?status.svg)](https://godoc.org/github.com/d4l3k/bypass)

bypass is an incredibly unsafe Go package that allows you to inspect private
parts of go objects.

* Fetches elements of a channel without modifying it
* Allows reflect to be used on unexported fields

## How to access private methods

https://sitano.github.io/2016/04/28/golang-private/

## License
Copyright (c) 2016 [Tristan Rice](https://fn.lc) <rice@fn.lc>

bypass is licensed under the MIT license. See the LICENSE file for more
information.

bypass.go and bypasssafe.go are borrowed from
[go-spew](https://github.com/davecgh/go-spew) and have a seperate copyright
notice.
