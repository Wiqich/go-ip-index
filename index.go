package ipindex

import (
	"fmt"
	"net"
)

const (
	// DefaultMinBinarySearchRange is the default value for minBinarySearchRange used for the situation the users won't decide this value by themself.
	DefaultMinBinarySearchRange = 10
)

// Value represents a additional value of the IP section. The values must be comparable.
type Value interface {
	Equal(other interface{}) bool
}

type section struct {
	lower uint32
	upper uint32
	value Value
}

func (sec section) Value() Value {
	return sec.value
}

// IPIndex represent a set of non-overlapping IP section and the addtional values.
// It use binary search and linear search to find the secion contains special IP and return its addtitional value.
type IPIndex struct {
	sections             []section
	index                [256][2]int
	minBinarySearchRange int
}

// Search the IP and return the additional value. It returns nil if the IP is not indexed.
// First, this function use index to narrow the search range by the A-Class IP section index.
// Second, this function use binary search to find the section contains the IP until it is found or the count of remain sections is smaller than minBinarySearchRange.
// Last, if it is not found in the second phase, this function use linear search to find the section contains the IP.
func (index *IPIndex) Search(ip net.IP) (Value, error) {
	key := ipToUint32(ip)
	if key == 0 {
		return nil, fmt.Errorf("invalid IP: %v", ip)
	}
	return index.search(key), nil
}

func (index *IPIndex) search(ip uint32) Value {
	initialRange := index.index[ip>>24]
	left, right := initialRange[0], initialRange[1]
	if left == -1 || right == -1 {
		return nil
	}
	for right-left > index.minBinarySearchRange {
		mid := (left + right) / 2
		if index.sections[mid].lower > ip {
			right = mid - 1
		} else if index.sections[mid].upper < ip {
			left = mid + 1
		} else {
			return index.sections[mid].value
		}
	}
	for i := left; i <= right; i++ {
		if index.sections[i].lower <= ip && index.sections[i].upper >= ip {
			return index.sections[i].value
		}
	}
	return nil
}
