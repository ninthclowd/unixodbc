package cache_test

import (
	"github.com/ninthclowd/unixodbc/internal/cache"
	"testing"
	"time"
)

func TestValue_Get(t *testing.T) {

	cachedIntValue := cache.Value[int]{}

	loadCount := 0

	loader := func() (int, error) {
		loadCount++
		return 5, nil
	}
	expiration := time.Now().Add(10 * time.Millisecond)

	got, err := cachedIntValue.Get(loader, &expiration)
	if got != 5 {
		t.Errorf("unexpected first value, got: %d", got)
	}
	if err != nil {
		t.Errorf("expected no error from first value, got: %s", err.Error())
	}
	if loadCount != 1 {
		t.Errorf("unexpected load count, got: %d", loadCount)
	}

	got2, err2 := cachedIntValue.Get(loader, &expiration)
	if got2 != 5 {
		t.Errorf("unexpected second value, got: %d", got2)
	}
	if err2 != nil {
		t.Errorf("expected no error from second value, got: %s", err2.Error())
	}

	time.Sleep(15 * time.Millisecond)

	if loadCount != 1 {
		t.Errorf("second call was not cached")
	}

	got3, err3 := cachedIntValue.Get(loader, &expiration)
	if got3 != 5 {
		t.Errorf("unexpected third value, got: %d", got3)
	}
	if err3 != nil {
		t.Errorf("expected no error from second value, got: %s", err3.Error())
	}
	if loadCount != 2 {
		t.Errorf("cache should expire with third call")
	}

}
