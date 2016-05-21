# go-ip-index
go-ip-index provides the IP index to index value associtated with non-overlapping IP sections.

[![Build Status](https://travis-ci.org/yangchenxing/go-ip-index.svg?branch=master)](https://travis-ci.org/yangchenxing/go-ip-index)
[![GoDoc](http://godoc.org/github.com/yangchenxing/go-ip-index?status.svg)](http://godoc.org/github.com/yangchenxing/go-ip-index)

##Example

    builder := NewIndexBuilder(DefaultMinBinarySearchRange)
    builder.Add(net.ParseIP("1.0.0.0"), net.ParserIP("1.0.0.255"), someValue)
    index := builder.Build()
    build.Search(net.ParseIP("1.0.0.1"))
