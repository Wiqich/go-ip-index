/*
go-ip-index provides the IP index to index value associtated with non-overlapping IP sections.

Example:

    builder := NewIndexBuilder(DefaultMinBinarySearchRange)
    builder.Add(net.ParseIP("1.0.0.0"), net.ParserIP("1.0.0.255"), someValue)
    index := builder.Build()
    build.Search(net.ParseIP("1.0.0.1"))
*/

package ipindex
