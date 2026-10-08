package pokecache

import (
	"fmt"
	"testing"
	"time"
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
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value %q, got %q", string(c.val), string(val))
				return
			}
		})
	}

	// Test getting non-existent key
	t.Run("non existent key", func(t *testing.T) {
		cache := NewCache(interval)
		_, ok := cache.Get("https://example.com/notfound")
		if ok {
			t.Errorf("expected not to find key")
		}
	})
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestReapFail(t *testing.T) {
	const baseTime = 50 * time.Millisecond
	const waitTime = 10 * time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	time.Sleep(waitTime)

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to still find key before interval expiration")
		return
	}
}

func BenchmarkGetParallel(b *testing.B) {
	cache := NewCache(time.Minute)
	cache.Add("https://example.com", []byte("testdata"))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cache.Get("https://example.com")
		}
	})
}
