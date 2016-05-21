package ipindex

import (
	"net"
	"testing"
)

func TestIndex(t *testing.T) {
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
			t.Errorf("builder.add fail: ip=%v, value=%v, error=%q\n", ips[i], values[i], err.Error())
			return
		}
	}
	index := builder.Build()
	value, err := index.Search(net.ParseIP("1.1.1.1"))
	if err != nil {
		t.Error("search 1.1.1.1 fail:", err.Error())
		return
	}
	if v, ok := value.(TestValue); !ok {
		t.Error("type assertion fail:", value)
		return
	} else if string(v) != "福州" {
		t.Errorf("unexpected search result: expected=%q, actual=%q\n", "福州", v)
	}
}
