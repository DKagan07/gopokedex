package pokecache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/DKagan07/gopokedex/pokecache"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second

	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const (
		baseTime = 50 * time.Millisecond
		waitTime = time.Second
		testUrl  = "https://example.com"
	)

	cache := pokecache.NewCache(baseTime)
	cache.Add(testUrl, []byte("testdata"))

	_, ok := cache.Get(testUrl)
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime * 3)

	_, ok = cache.Get(testUrl)
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
