package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cases := []struct {
		key   string
		val   []byte
		found bool
	}{
		{
			key: "first",
			val: []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		},
		{
			key: "second",
			val: []byte("second test beep boop"),
		},
	}
	for _, c := range cases {
		csh := NewCache(5 * time.Second)
		csh.Add(c.key, c.val)
		val, found := csh.Get(c.key)
		if !found {
			t.Errorf("expected to find key")
			return
		}
		if string(val) != string(c.val) {
			t.Errorf("expected data to remain unchanged\n%s : %s", string(val), string(c.val))
		}
	}
}

func TestNoAddGet(t *testing.T) {
	cases := []struct {
		key   string
		val   []byte
		found bool
	}{
		{
			key: "first",
		},
		{
			key: "second",
		},
	}
	for _, c := range cases {
		csh := NewCache(5 * time.Second)

		val, found := csh.Get(c.key)
		if found {
			t.Errorf("expected to not find data")
		} else if val != nil {
			t.Errorf("expected data to be nil")
		}
	}
}

func TestReapLoop(t *testing.T) {

	testUrl := "https://example.com"
	testVal := []byte("foobar")

	baseTime := 5 * time.Millisecond
	waitTime := baseTime + (5 * time.Millisecond)
	csh := NewCache(baseTime)

	csh.Add(testUrl, testVal)

	_, found := csh.Get(testUrl)

	if !found {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, found = csh.Get(testUrl)
	if found {
		t.Errorf("expected to not find key")
		return
	}
}
