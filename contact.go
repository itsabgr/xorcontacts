package xorcontacts

import (
	"bytes"
	"time"
)

type Contact interface {
	ID() []byte
}
type Bytes []byte

func (b Bytes) ID() []byte {
	return b
}

type contact[C Contact] struct {
	wrapped  C
	deadline time.Time
}

func (c *contact[C]) Compare(another *contact[C]) int {
	return bytes.Compare(c.wrapped.ID(), another.wrapped.ID())
}
func (c *contact[C]) Expired() bool {
	return !c.deadline.After(time.Now())
}
