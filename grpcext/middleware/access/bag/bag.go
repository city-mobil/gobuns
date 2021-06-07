package bag

import (
	"context"
)

// key is a type for keys defined in this package.
type key int

// bagKey is the key for bag.Bag values in context.
var bagKey key

type Field struct {
	Key   string
	Value interface{}
}

// Bag carries user fields to log them in the access interceptor.
type Bag map[string]interface{}

func (b Bag) Add(fields ...Field) {
	if b == nil {
		return
	}

	for _, f := range fields {
		b[f.Key] = f.Value
	}
}

func New() Bag {
	return make(Bag)
}

func NewContext(ctx context.Context, b Bag) context.Context {
	return context.WithValue(ctx, bagKey, b)
}

func FromContext(ctx context.Context) (Bag, bool) {
	b, ok := ctx.Value(bagKey).(Bag)
	return b, ok
}
