package ipindex

import (
	"net"
	"testing"
)

func TestIndex(t *testing.T) {
	builder := NewIndexBuilder(DefaultMinBinarySearchRange)
	ips := [][2]string{
		{"1.0.0.0", "1.0.0.255"},
		{"1.0.1.0", "1.0.1.255"},
		{"1.0.2.0", "1.0.2.255"},
		{"1.0.3.0", "1.0.3.255"},
		{"1.0.4.0", "3.0.4.255"},
		{"3.0.6.0", "3.0.6.255"},
		{"3.0.7.0", "3.0.7.255"},
		{"3.0.8.0", "3.0.8.255"},
		{"3.0.9.0", "3.0.9.255"},
		{"3.0.10.0", "3.0.10.255"},
		{"3.0.11.0", "3.0.11.255"},
		{"3.0.12.0", "3.0.12.255"},
		{"3.0.13.0", "3.0.13.255"},
		{"3.0.14.0", "3.0.14.255"},
		{"3.0.15.0", "3.0.15.255"},
		{"3.0.16.0", "3.0.16.255"},
		{"3.0.17.0", "3.0.17.255"},
		{"3.0.18.0", "3.0.18.255"},
		{"3.0.19.0", "3.0.19.255"},
		{"3.0.20.0", "3.0.20.255"},
		{"3.0.21.0", "3.0.21.255"},
		{"3.0.22.0", "3.0.22.255"},
		{"3.0.23.0", "3.0.23.255"},
		{"3.0.24.0", "3.0.24.255"},
		{"3.0.25.0", "3.0.25.255"},
		{"3.0.26.0", "3.0.26.255"},
		{"3.0.27.0", "3.0.27.255"},
		{"3.0.28.0", "3.0.28.255"},
		{"3.0.29.0", "3.0.29.255"},
		{"3.0.30.0", "3.0.30.255"},
		{"3.0.31.0", "3.0.31.255"},
		{"3.0.32.0", "3.0.32.255"},
		{"3.0.33.0", "3.0.33.255"},
		{"3.0.34.0", "3.0.34.255"},
		{"3.0.35.0", "3.0.35.255"},
		{"3.0.36.0", "3.0.36.255"},
		{"3.0.37.0", "3.0.37.255"},
		{"3.0.38.0", "3.0.38.255"},
		{"3.0.39.0", "3.0.39.255"},
		{"3.0.40.0", "3.0.40.255"},
		{"3.0.41.0", "3.0.41.255"},
		{"3.0.42.0", "3.0.42.255"},
		{"3.0.43.0", "3.0.43.255"},
		{"3.0.44.0", "3.0.44.255"},
		{"3.0.45.0", "3.0.45.255"},
		{"3.0.46.0", "3.0.46.255"},
		{"3.0.47.0", "3.0.47.255"},
		{"3.0.48.0", "3.0.48.255"},
		{"3.0.49.0", "3.0.49.255"},
		{"3.0.50.0", "3.0.50.255"},
		{"3.0.51.0", "3.0.51.255"},
		{"3.0.52.0", "3.0.52.255"},
		{"3.0.53.0", "3.0.53.255"},
		{"3.0.54.0", "3.0.54.255"},
	}
	values := []TestValue{
		TestValue("上海"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("福州"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
		TestValue("福州"),
		TestValue("上海"),
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
	// invalid ip
	if _, err := index.Search(nil); err == nil {
		t.Error("unexpected search success")
		return
	}
	// search index miss entities
	if value, err := index.Search(net.ParseIP("127.0.0.1")); err != nil {
		t.Error("search fail:", err.Error())
		return
	} else if value != nil {
		t.Error("unexpected search result:", value)
		return
	}
	// binary search
	if value, err := index.Search(net.ParseIP("3.0.22.100")); err != nil {
		t.Error("search fail:", err.Error())
		return
	} else if !value.Equal(TestValue("上海")) {
		t.Error("unexpected search result:", value)
		return
	}
	// search miss
	if value, err := index.Search(net.ParseIP("3.1.0.100")); err != nil {
		t.Error("search fail:", err.Error())
		return
	} else if value != nil {
		t.Error("unexpected search result:", value)
		return
	}
}
