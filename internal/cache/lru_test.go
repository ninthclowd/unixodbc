package cache_test

import (
	"errors"
	"fmt"
	"github.com/ninthclowd/unixodbc/internal/cache"
	"strconv"
	"testing"
)

type MyStruct struct{}

func TestLRU_Get_Empty(t *testing.T) {
	lru := cache.NewLRU[MyStruct](0, nil)

	if gotValue := lru.Get("A", false); gotValue != nil {
		t.Fatalf("expected nil but received: %v", gotValue)
	}
}

func TestLRU_Put_Empty(t *testing.T) {
	lru := cache.NewLRU[MyStruct](0, nil)
	if gotErr := lru.Put("A", new(MyStruct)); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}
}

func TestLRU_Put_Empty_Eviction(t *testing.T) {
	var wantValue MyStruct
	wantKey := "A"
	wantErr := errors.New("foo")

	onEvict := func(gotKey string, gotValue *MyStruct) error {
		if gotKey != wantKey {
			t.Fatalf("expected key of %s but got %s", wantKey, gotKey)
		}
		if gotValue != &wantValue {
			t.Fatalf("expected value of %v but got %v", wantValue, gotValue)
		}
		return wantErr
	}

	lru := cache.NewLRU[MyStruct](0, onEvict)

	if gotErr := lru.Put("A", &wantValue); !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected error %v but received: %v", wantErr, gotErr)
	}
}

func TestLRU_Put_Evict_Error(t *testing.T) {
	var wantValue MyStruct
	wantKey := "A"
	wantErr := errors.New("foo")

	onEvict := func(gotKey string, gotValue *MyStruct) error {
		if gotKey != wantKey {
			t.Fatalf("expected key of %s but got %s", wantKey, gotKey)
		}
		if gotValue != &wantValue {
			t.Fatalf("expected value of %v but got %v", wantValue, gotValue)
		}
		return wantErr
	}

	lru := cache.NewLRU[MyStruct](1, onEvict)

	if gotErr := lru.Put(wantKey, &wantValue); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}

	if gotErr := lru.Put("B", new(MyStruct)); !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected error %v but received %v", wantErr, gotErr)
	}
}

func TestLRU_Put_Redundant(t *testing.T) {
	var wantValue MyStruct

	lru := cache.NewLRU[MyStruct](1, nil)

	if gotErr := lru.Put("A", &wantValue); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}

	if gotErr := lru.Put("A", &wantValue); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}
}

func TestLRU_Purge(t *testing.T) {
	var itemA, itemB MyStruct
	evictions := make(map[string]*MyStruct)
	onEvict := func(key string, value *MyStruct) error {
		evictions[key] = value
		return nil
	}
	lru := cache.NewLRU[MyStruct](2, onEvict)

	if gotErr := lru.Put("A", &itemA); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}
	if gotErr := lru.Put("B", &itemB); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}
	if gotEvictions := len(evictions); gotEvictions != 0 {
		t.Fatalf("expected no evictions but received: %d", gotEvictions)
	}

	if gotErr := lru.Purge(); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}

	if gotEvictions := len(evictions); gotEvictions != 2 {
		t.Fatalf("expected 2 evictions but received: %d", gotEvictions)
	}
}

func TestLRU_Purge_Error(t *testing.T) {
	wantErr := errors.New("foo")
	onEvict := func(key string, value *MyStruct) error {
		return wantErr
	}
	lru := cache.NewLRU[MyStruct](1, onEvict)

	if gotErr := lru.Put("A", new(MyStruct)); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}

	if gotErr := lru.Purge(); !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected %v but received %v", wantErr, gotErr)
	}
}

func TestLRU_Get_Remove(t *testing.T) {
	var itemA MyStruct
	lru := cache.NewLRU[MyStruct](1, func(key string, value *MyStruct) error {
		t.Fatalf("onEvict should not be called when value is removed by Get")
		return nil
	})
	if gotErr := lru.Put("A", &itemA); gotErr != nil {
		t.Fatalf("expected no error but received: %v", gotErr)
	}

	if gotValue := lru.Get("A", true); gotValue != &itemA {
		t.Fatalf("expected %v error but received: %v", &itemA, gotValue)
	}

	if gotValue := lru.Get("A", true); gotValue != nil {
		t.Fatal("expected item to be remove from previous call to Get")
	}

}

func TestLRU_Get(t *testing.T) {
	var itemA, itemB, itemC, itemD MyStruct

	evictions := make(map[string]*MyStruct)
	onEvict := func(key string, value *MyStruct) error {
		evictions[key] = value
		return nil
	}
	lru := cache.NewLRU[MyStruct](2, onEvict)

	type Test struct {
		putKey      string
		putVal      *MyStruct
		wantEvicted map[string]*MyStruct
		wantGet     map[string]*MyStruct
		wantLen     int
	}

	tests := []Test{
		{
			putKey: "A",
			putVal: &itemA,
			wantGet: map[string]*MyStruct{
				"A": &itemA,
			},
			wantEvicted: nil,
			wantLen:     1,
		},
		{
			putKey: "B",
			putVal: &itemB,
			wantGet: map[string]*MyStruct{
				"A": &itemA,
				"B": &itemB,
			},
			wantEvicted: nil,
			wantLen:     2,
		},
		{
			putKey: "C",
			putVal: &itemC,
			wantGet: map[string]*MyStruct{
				"B": &itemB,
				"C": &itemC,
			},
			wantEvicted: map[string]*MyStruct{
				"A": &itemA,
			},
			wantLen: 2,
		},
		{
			putKey: "D",
			putVal: &itemD,
			wantGet: map[string]*MyStruct{
				"C": &itemC,
				"D": &itemD,
			},
			wantEvicted: map[string]*MyStruct{
				"A": &itemA,
				"B": &itemB,
			},
			wantLen: 2,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("After adding %s", test.putKey), func(t *testing.T) {

			if err := lru.Put(test.putKey, test.putVal); err != nil {
				t.Fatalf("expected no error from Put.  Got %v", err)
			}

			for key, value := range test.wantGet {
				if gotValue := lru.Get(key, false); gotValue != value {
					t.Fatalf("item %s was expected to be loaded but was not", key)
				}
			}

			for key, value := range test.wantEvicted {
				if gotValue := lru.Get(key, false); gotValue != nil {
					t.Fatalf("item %s was expected to be evicted but was not", key)
				}
				if evictedValue, evicted := evictions[key]; !evicted {
					t.Fatalf("expected onEvict to be called with putKey %s", key)
				} else if evictedValue != value {
					t.Fatalf("putVal sent to onEvict for putKey %v was unexpected", key)
				}

			}
		})
	}
}

func BenchmarkLRU_Put(b *testing.B) {
	capacity := 5

	lru := cache.NewLRU[MyStruct](capacity, nil)

	//prefill the cache to capacity
	i := 0
	for i < capacity {
		_ = lru.Put(strconv.Itoa(i), new(MyStruct))
		i++
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = lru.Put(strconv.Itoa(i), new(MyStruct))
		i++
	}
	b.StopTimer()
}
