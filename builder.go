package ipindex

import (
	"container/list"
	"errors"
	"fmt"
	"net"
)

// IndexBuilder build the IPIndex by add non-overlapping ip sections and the additional values in ascending order.
type IndexBuilder struct {
	sections             *list.List
	last                 section
	minBinarySearchRange int
}

// NewIndexBuilder create a new IndexBuilder instance and specified the minBinarySearchRange parameter for the IPIndex it will build.
func NewIndexBuilder(minBinarySearchRange int) *IndexBuilder {
	return &IndexBuilder{
		sections:             list.New(),
		minBinarySearchRange: minBinarySearchRange,
	}
}

// Reset the builder to initial state.
func (builder *IndexBuilder) Reset() {
	builder.sections.Init()
}

// Add a IP section and its addtional value. The lower bound must be bigger than last section added.
// If the lower bound is adjacent to the upper bound of last section, and their are same kind of A-Class address,
// and the value is equal to that of last section, the new section will merge to last section.
func (builder *IndexBuilder) Add(lower, upper net.IP, value Value) error {
	if lower := ipToUint32(lower); lower == 0 {
		return fmt.Errorf("invalid lower: %v", lower)
	} else if upper := ipToUint32(upper); upper == 0 {
		return fmt.Errorf("invalid upper: %v", upper)
	} else {
		return builder.add(lower, upper, value)
	}
}

func (builder *IndexBuilder) add(lower, upper uint32, value Value) error {
	// 检查参数
	if lower == 0 {
		return errors.New("lower bound cannot be 0")
	}
	if lower > upper {
		return fmt.Errorf("lower bound is higher than upper bound: lower=%d, upper=%d",
			lower, upper)
	}
	// 切分跨A类长区间
	if lower>>24 < upper>>24 {
		aLower := lower >> 24
		aUpper := upper >> 24
		if err := builder.add(lower, (aLower<<24)|0x00FFFFFF, value); err != nil {
			return err
		}
		for i := aLower + 1; i < aUpper; i++ {
			if err := builder.add(i<<24, (i<<24)|0x00FFFFFF, value); err != nil {
				return err
			}
		}
		if err := builder.add(aUpper<<24, upper, value); err != nil {
			return err
		}
		return nil
	}
	// 赋值首个区间
	if builder.last.lower == 0 {
		builder.last.lower = lower
		builder.last.upper = upper
		builder.last.value = value
		return nil
	}
	// 合并相邻同A段同Value的区间
	if builder.last.upper+1 == lower && builder.last.value.Equal(value) &&
		builder.last.upper>>24 == lower>>24 {
		builder.last.upper = upper
		return nil
	}
	// 保存
	builder.sections.PushBack(builder.last)
	builder.last.lower = lower
	builder.last.upper = upper
	builder.last.value = value
	return nil
}

// Build the IPIndex.
func (builder *IndexBuilder) Build() *IPIndex {
	// 保存最后一个区间
	if builder.last.lower > 0 {
		builder.sections.PushBack(builder.last)
		builder.last.lower = 0
		builder.last.upper = 0
		builder.last.value = nil
	}
	sections := make([]section, builder.sections.Len())
	var index [256][2]int
	for i := 0; i < 256; i++ {
		index[i][0] = -1
		index[i][1] = -1
	}
	for i, sec := 0, builder.sections.Front(); sec != nil && i < len(sections); i, sec = i+1, sec.Next() {
		section := sec.Value.(section)
		sections[i] = section
		indexKey := section.lower >> 24
		if index[indexKey][0] == -1 {
			index[indexKey][0] = i
		}
		index[indexKey][1] = i
	}
	return &IPIndex{
		sections:             sections,
		index:                index,
		minBinarySearchRange: builder.minBinarySearchRange,
	}
}

func ipToUint32(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
}
