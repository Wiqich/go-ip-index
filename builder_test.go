package ipindex

import (
	"net"
	"testing"
)

type TestValue string

func (v TestValue) Equal(other interface{}) bool {
	w, ok := other.(TestValue)
	if !ok {
		return false
	}
	return w == v
}

func TestBuilder(t *testing.T) {
	builder := NewIndexBuilder(DefaultMinBinarySearchRange)
	ips := [][2]string{
		[2]string{"1.0.0.0", "1.0.0.255"},
		[2]string{"1.0.1.0", "1.0.1.255"},
		[2]string{"1.0.2.0", "1.0.2.255"},
		[2]string{"1.0.3.0", "1.0.3.255"},
		[2]string{"1.0.4.0", "2.0.4.255"},
	}
	values := []TestValue{
		TestValue("上海"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("福州"),
		TestValue("福州"),
	}
	for i := 0; i < len(ips); i++ {
		if err := builder.Add(net.ParseIP(ips[i][0]), net.ParseIP(ips[i][1]), values[i]); err != nil {
			t.Error("builder.add fail: ip=%v, value=%v, error=%q", ips[i], values[i], err.Error())
			return
		}
	}
	index := builder.Build()
	expectedSections := []section{
		{
			lower: 0x01000000,
			upper: 0x010001FF,
			value: TestValue("上海"),
		},
		{
			lower: 0x01000200,
			upper: 0x01FFFFFF,
			value: TestValue("福州"),
		},
		{
			lower: 0x02000000,
			upper: 0x020004FF,
			value: TestValue("福州"),
		},
	}
	if len(expectedSections) != len(index.sections) {
		t.Error("unexpected sections: expected=%v, actual=%v", expectedSections, index.sections)
		return
	}
	for i, exp := range expectedSections {
		act := index.sections[i]
		if exp.lower != act.lower || exp.upper != act.upper || !exp.Value().Equal(act.Value()) {
			t.Error("unexpected sections: expected=%v, actual=%v", expectedSections, index.sections)
			return
		}
	}
	var expectedIndex [256][2]int
	for i := 0; i < 256; i++ {
		expectedIndex[i][0] = -1
		expectedIndex[i][1] = -1
	}
	expectedIndex[1][0] = 0
	expectedIndex[1][1] = 1
	expectedIndex[2][0] = 2
	expectedIndex[2][1] = 2
	for i := 0; i < 256; i++ {
		if index.index[i][0] != expectedIndex[i][0] || index.index[i][1] != expectedIndex[i][1] {
			t.Error("unexpected index[%d]: expected=%v, actual=%v", i, expectedIndex[i], index.index[i])
			return
		}
	}
}
