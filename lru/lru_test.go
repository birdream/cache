package lru

import (
	"reflect"
	"testing"
)

type Str string

func (s Str) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := NewCache(int64(0), nil)
	lru.Add("key1", Str("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(Str)) != "1234" {
		t.Fatalf("get hit key1=1234 fail")
	}

	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("get should not return key2")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "val1", "val2", "val3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap), nil)

	lru.Add(k1, Str(v1))
	lru.Add(k2, Str(v2))
	lru.Add(k3, Str(v3))

	if _, ok := lru.Get(k1); ok || lru.Len() != 2 {
		t.Fatalf("RemoveOldesr not work as expect")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	cb := func(key string, val Value) {
		keys = append(keys, key)
	}

	lru := NewCache(int64(10), cb)

	lru.Add("key1", Str("123456"))
	lru.Add("k2", Str("k2"))
	lru.Add("k3", Str("k3"))
	lru.Add("k4", Str("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call onEvicted failes")
	}
}
