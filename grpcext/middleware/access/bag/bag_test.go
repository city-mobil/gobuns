package bag

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBag_Add(t *testing.T) {
	b := New()
	b.Add(Field{"str", "b"}, Field{"int", 1})
	assert.Len(t, b, 2)
}

func TestBag_Add_Nil(t *testing.T) {
	var b Bag
	b.Add(Field{"a", "b"})
}

func TestBag_Context(t *testing.T) {
	b := New()
	ctx := NewContext(context.Background(), b)
	got, ok := FromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, b, got)
}
