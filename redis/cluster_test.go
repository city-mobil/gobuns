package redis

import (
	"context"
	"testing"
	"time"

	"github.com/city-mobil/gobuns/barber"
	"github.com/city-mobil/gobuns/zlog"
)

func TestNewCluster(t *testing.T) {
	c := newDefaultCluster(
		zlog.New(nil),
		newLocalStorage(),
		&fallback{},
		barber.NewBarber([]int{0}, &barber.Config{
			Threshold: 100,
			MaxFails:  100,
		}),
	)

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
