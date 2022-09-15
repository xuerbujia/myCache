package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(bytes []byte) uint32 {
		i, _ := strconv.Atoi(string(bytes))
		return uint32(i)
	})
	hash.Add("2", "4", "6")
	testCase := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
		"15": "6",
	}
	for k, v := range testCase {
		get := hash.Get(k)
		if get != v {
			t.Fatalf("get%v,excepted%v", get, v)
		}
	}
	hash.Add("8")
	testCase["27"] = "8"
	for k, v := range testCase {
		get := hash.Get(k)
		if get != v {
			t.Fatalf("k:%v,get%v,excepted%v", k, get, v)
		}
	}
}
