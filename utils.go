package xorcontacts

import xorLib "github.com/itsabgr/xorcontacts/pkg/xor"

func xor(a, b []byte) []byte {
	dst := make([]byte, len(a))
	xorLib.XorBytes(dst, a, b)
	return dst
}
