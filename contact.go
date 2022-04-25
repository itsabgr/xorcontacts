package xorcontacts

import (
	"bytes"
	"time"
)

type Contact interface {
	UUID() []byte
	Hash() []byte
}

type contact[C Contact] struct {
	wrapped  C
	deadline time.Time
}

func (c *contact[C]) Compare(another *contact[C]) int {
	return bytes.Compare(c.wrapped.UUID(), another.wrapped.UUID())
}
func (c *contact[C]) Expired() bool {
	return !c.deadline.After(time.Now())
}
