package zlog

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	log := New(ioutil.Discard)
	ctx := NewContext(context.Background(), log)
	log2 := FromContext(ctx)
	assert.Equal(t, log, log2)

	log = log.Level(InfoLevel)
	ctx = NewContext(ctx, log)
	log2 = FromContext(ctx)
	assert.Equal(t, log, log2)

	log2 = FromContext(context.Background())
	assert.Same(t, disabledLogger, log2)
}

func TestContextOverride(t *testing.T) {
	l := New(ioutil.Discard).With().Str("foo", "bar").Logger()
	ctx := NewContext(context.Background(), l)
	assert.Same(t, l, FromContext(ctx), "NewContext did not store logger")

	l.UpdateContext(func(c Context) Context {
		return c.Str("bar", "baz")
	})
	ctx = NewContext(ctx, l)
	assert.Same(t, l, FromContext(ctx), "NewContext did not store updated logger")

	l = l.Level(DebugLevel)
	ctx = NewContext(ctx, l)
	assert.Same(t, l, FromContext(ctx), "NewContext did not store copied logger")
}

func TestContextDisabled(t *testing.T) {
	dl := New(ioutil.Discard).Level(Disabled)
	ctx := NewContext(context.Background(), dl)
	assert.Equal(t, context.Background(), ctx, "NewContext stored a disabled logger")

	ctx = NewContext(ctx, dl)
	assert.Equal(t, dl, FromContext(ctx), "NewContext did not override logger with a disabled logger")
}
