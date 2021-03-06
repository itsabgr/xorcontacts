package xorcontacts

import (
	"crypto/rand"
	"github.com/itsabgr/go-handy"
	"io"
	"sort"
	"testing"
	"time"
)

func RandomPeer() Bytes {
	b := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, b)
	handy.Throw(err)
	return b
}

func TestList(t *testing.T) {
	list := List[Bytes]{}
	for range handy.N(100) {
		if nil != list.Add(RandomPeer(), time.Now().Add(1*time.Hour)) {
			t.Fatal()
		}
	}
	sort.Sort(list.xor(RandomPeer()))
}
func TestList_GC(t *testing.T) {
	list := List[Bytes]{}
	dl := time.Now().Add(3 * time.Second)
	for range handy.N(100) {
		if nil != list.Add(RandomPeer(), dl) {
			t.Fatal()
		}
	}
	dl = time.Now().Add(20 * time.Second)
	for range handy.N(100) {
		if nil != list.Add(RandomPeer(), dl) {
			t.Fatal()
		}
	}
	if list.Len() != 200 {
		t.Fatal()
	}
	<-time.NewTimer(4 * time.Second).C
	list.GC()
	if list.Len() != 100 {
		t.Fatal(list.Len())
	}
}
