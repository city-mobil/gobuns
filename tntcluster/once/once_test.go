package once

import "testing"

func TestOnceReset(t *testing.T) {
	var calls int
	var c Once
	c.Do(func() {
		calls++
	})
	c.Do(func() {
		calls++
	})
	c.Do(func() {
		calls++
	})
	if calls != 1 {
		t.Errorf("Call count equal to %v expected 1", calls)
	}
	c.Reset()
	c.Do(func() {
		calls++
	})
	c.Do(func() {
		calls++
	})
	c.Do(func() {
		calls++
	})
	if calls != 2 {
		t.Errorf("Call count equal to %v expected 2", calls)
	}
}
