# ifacecodegen

[![Build Status](https://travis-ci.org/fredipevcin/ifacecodegen.svg?branch=master)](https://travis-ci.org/fredipevcin/ifacecodegen)
[![GoDoc](https://godoc.org/github.com/fredipevcin/ifacecodegen?status.svg)](https://godoc.org/github.com/fredipevcin/ifacecodegen)

Go tool and library for generating code from the template using interface definition.

## Installation

	go get -u github.com/fredipevcin/ifacecodegen/cmd/ifacecodegen


## Running ifacecodegen

	ifacecodegen -source examples/interface.go -destination -

or

	cat examples/interface.go | ifacecodegen

Other options

    ifacecodegen -h

## Examples

Example templates and interface are located in `examples` folder.

```bash
# example 1
ifacecodegen \
	-source examples/interface.go \
	-template examples/example1.tmpl \
	-destination - \
	-meta service=account \
	-imports "opentracing=github.com/opentracing/opentracing-go,tracinglog=github.com/opentracing/opentracing-go/log"

# example 2
ifacecodegen \
	-source examples/interface.go \
	-template examples/example2.tmpl \
	-destination -
```


## Notes

Inspired by:

* https://github.com/golang/mock
* https://github.com/kevinconway/wrapgen
