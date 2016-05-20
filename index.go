package ipindex

import (
	"fmt"
	"net"
)

const (
	DefaultMinBinarySearchRange = 10
)

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

type IPIndex struct {
	sections             []section
	index                [256][2]int
	minBinarySearchRange int
}

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
