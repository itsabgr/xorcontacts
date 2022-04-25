package xorcontacts

import (
	"crypto/md5"
	xorLib "github.com/itsabgr/xorcontacts/pkg/xor"
)

func xor(a, b []byte) []byte {
	dst := make([]byte, len(a))
	xorLib.XorBytes(dst, a, b)
	return dst
}

func hash(b []byte) []byte {
	h := md5.Sum(b)
	return h[:]
}
