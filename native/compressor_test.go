package native

import "testing"

var shortString = []byte("hello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\n")

func TestNewCompressor(t *testing.T) {
	c, err := NewCompressor(defaultLevel)
	defer c.Close()
	if err != nil {
		t.Error(err)
	}

	c, err = NewCompressor(30)
	if err == nil {
		t.Fail()
	}
}

// this test doesn't really say much - integration tests will show correctnes
func TestCompress(t *testing.T) {
	c, _ := NewCompressor(defaultLevel)
	defer c.Close()

	if _, _, err := c.Compress(make([]byte, 0), nil); err == nil {
		t.Error("expected error")
	}

	n, out, err := c.Compress(shortString, nil)
	if err != nil || n == 0 || n >= len(shortString) || n != len(out) {
		t.Error(err)
		t.Error(n)
	}

	out2 := make([]byte, len(shortString))
	n, _, err = c.Compress(shortString, out2)
	if err != nil || n == 0 {
		t.Error(err)
		t.Error(n)
	}

	slicesEqual(out, out2[:n], t) //tests if rep produces same results
}

func slicesEqual(expected, actual []byte, t *testing.T) {
	if len(expected) != len(actual) {
		t.Error("len of slices unequal")
		t.FailNow()
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("slices differ at %d: want: %d, got: %d", i, expected[i], actual[i])
			t.FailNow()
		}
	}
}
