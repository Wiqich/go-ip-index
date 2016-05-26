# go-ip-index
go-ip-index provides the IP index to index value associtated with non-overlapping IP sections.

[![Go Report Card](https://goreportcard.com/badge/github.com/yangchenxing/go-ip-index)](https://goreportcard.com/report/github.com/yangchenxing/go-ip-index)
[![Build Status](https://travis-ci.org/yangchenxing/go-ip-index.svg?branch=master)](https://travis-ci.org/yangchenxing/go-ip-index)
[![GoDoc](http://godoc.org/github.com/yangchenxing/go-ip-index?status.svg)](http://godoc.org/github.com/yangchenxing/go-ip-index)
[![Coverage Status](https://coveralls.io/repos/github/yangchenxing/go-ip-index/badge.svg?branch=master)](https://coveralls.io/github/yangchenxing/go-ip-index?branch=master)

##Example

    builder := NewIndexBuilder(DefaultMinBinarySearchRange)
    builder.Add(net.ParseIP("1.0.0.0"), net.ParserIP("1.0.0.255"), someValue)
    index := builder.Build()
    build.Search(net.ParseIP("1.0.0.1"))
