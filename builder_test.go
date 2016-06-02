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
		[2]string{"1.0.4.0", "3.0.4.255"},
		[2]string{"3.0.6.0", "3.0.6.255"},
	}
	values := []TestValue{
		TestValue("上海"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("福州"),
		TestValue("福州"),
		TestValue("上海"),
	}
	for i := 0; i < len(ips); i++ {
		if err := builder.Add(net.ParseIP(ips[i][0]), net.ParseIP(ips[i][1]), values[i]); err != nil {
			t.Errorf("builder.add fail: ip=%v, value=%v, error=%q\n", ips[i], values[i], err.Error())
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
			upper: 0x02FFFFFF,
			value: TestValue("福州"),
		},
		{
			lower: 0x03000000,
			upper: 0x030004FF,
			value: TestValue("福州"),
		},
		{
			lower: 0x03000600,
			upper: 0x030006FF,
			value: TestValue("上海"),
		},
	}
	if len(expectedSections) != len(index.sections) {
		t.Errorf("unexpected sections: expected=%v, actual=%v\n", expectedSections, index.sections)
		return
	}
	for i, exp := range expectedSections {
		act := index.sections[i]
		if exp.lower != act.lower || exp.upper != act.upper || !exp.Value().Equal(act.Value()) {
			t.Errorf("unexpected sections: expected=%v, actual=%v\n", expectedSections, index.sections)
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
	expectedIndex[3][0] = 3
	expectedIndex[3][1] = 4
	for i := 0; i < 256; i++ {
		if index.index[i][0] != expectedIndex[i][0] || index.index[i][1] != expectedIndex[i][1] {
			t.Errorf("unexpected index[%d]: expected=%v, actual=%v\n", i, expectedIndex[i], index.index[i])
			return
		}
	}
	// reset
	builder.Reset()
	if builder.sections.Len() > 0 {
		t.Error("reset fail")
		return
	}
	// lower == 0
	if err := builder.Add(net.ParseIP("0.0.0.0"), net.ParseIP("0.0.0.1"), TestValue("福州")); err == nil {
		t.Error("unexpected add success: lower == 0")
		return
	}
	if err := builder.AddUint32(0, 1, TestValue("福州")); err == nil {
		t.Error("unexpected add success: lower == 0")
		return
	}
	// upper == 0
	if err := builder.Add(net.ParseIP("0.0.0.1"), net.ParseIP("0.0.0.0"), TestValue("福州")); err == nil {
		t.Error("unexpected add success: upper == 0")
		return
	}
	if err := builder.AddUint32(1, 0, TestValue("福州")); err == nil {
		t.Error("unexpected add success: upper == 0")
		return
	}
	// upper < lower
	if err := builder.AddUint32(100, 1, TestValue("福州")); err == nil {
		t.Error("unexpected add success: lower > upper")
		return
	}
}

func TestIPToUint32(t *testing.T) {
	if value := ipToUint32(nil); value != 0 {
		t.Error("unexpected result:", value)
		return
	}
	if value := ipToUint32(net.ParseIP("FE80:0000:0000:0000:0202:B3FF:FE1E:8329")); value != 0 {
		t.Error("unexpected result:", value)
		return
	}
}
