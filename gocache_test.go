package gocache

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadcounter := make(map[string]int, len(db))
	g := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			t.Log("[slow query]search key", key)
			if v, ok := db[key]; ok {
				if c, ok := loadcounter[key]; !ok {
					loadcounter[key] = 1
				} else {
					loadcounter[key] = c + 1
				}
				return []byte(v), nil
			}
			return nil, fmt.Errorf("key %s not exists", key)
		}))
	for k, v := range db {
		bv, err := g.Get(k)
		if err != nil {
			t.Fatalf("group get error")
		}
		sc := bv.String()
		if sc != db[k] {
			t.Fatalf("key %s value expected %s but get %s", k, v, sc)
		}
		if loadcounter[k] != 1 {
			t.Fatalf("cache missed")
		}
	}

	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
