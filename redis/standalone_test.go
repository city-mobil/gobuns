package redis

import (
	"context"
	"testing"
	"time"

	"github.com/city-mobil/gobuns/barber"
)

func TestNewStandalone(t *testing.T) {
	repls := make([]*node, 0)
	repls = append(repls, &node{
		newLocalStorage(),
		1,
	})
	cb := barber.NewBarber([]int{0, 1}, &barber.Config{
		Threshold: 100,
		MaxFails:  100,
	})

	conn := newStandaloneConn(
		&node{
			item:  newLocalStorage(),
			index: 0,
		},
		repls,
		&node{
			item:  &fallback{},
			index: -1,
		},
		cb,
	)
	c := newDefaultStandalone(conn)

	ctx := context.Background()

	_, _ = c.Set(ctx, "hello", "world", time.Hour)
	val, err := c.Get(ctx, "hello")
	if err != nil {
		t.Errorf("didn't expected error %s", err.Error())
	}

	if val != "world" {
		t.Errorf("can not sotre value. expected: world got: %s", val)
	}
}
