package xorcontacts

import (
	"bytes"
	"errors"
	"github.com/itsabgr/go-set"
	"time"
)

type List[C Contact] struct {
	set set.Set[*contact[C]]
}

func (l *List[C]) Add(peer C, deadline time.Time) error {
	if false == deadline.After(time.Now()) {
		return errors.New("expired deadline")
	}
	if false == l.set.Add(&contact[C]{deadline: deadline, wrapped: peer}) {
		return errors.New("exists")
	}
	return nil
}
func (l *List[C]) Update(peer C, deadline time.Time) error {
	if false == deadline.After(time.Now()) {
		return errors.New("expired deadline")
	}
	if !l.set.Remove(&contact[C]{wrapped: peer}) {
		return errors.New("not found")
	}
	return l.Add(peer, deadline)
}
func (l *List[C]) Upsert(peer C, deadline time.Time) error {
	if false == deadline.After(time.Now()) {
		return errors.New("expired deadline")
	}
	_ = l.set.Remove(&contact[C]{wrapped: peer})
	return l.Add(peer, deadline)
}
func (l *List[C]) Remove(peer C) error {
	if false == l.set.Remove(&contact[C]{wrapped: peer}) {
		return errors.New("not found")
	}
	return nil
}
func (l *List[C]) GC() {
	var todo []*contact[C]
	for iter := l.set.Iter(); iter.HasNext(); iter.Next() {
		if false == iter.Item().deadline.After(time.Now()) {
			todo = append(todo, iter.Item())
		}
	}
	for _, item := range todo {
		l.set.Remove(item)
	}
}

func (l *List[C]) Has(peer C) bool {
	return l.set.Has(&contact[C]{wrapped: peer})
}
func (l *List[C]) Len() int {
	return l.set.Len()
}
func (l *List[C]) List() set.Set[*contact[C]] {
	return l.set
}

type PeerWithXor[C Contact] struct {
	peer C
	xor  []byte
}

func (p *PeerWithXor[C]) Peer() C {
	return p.peer
}

func (p *PeerWithXor[C]) Xor() []byte {
	return p.xor
}

type ListPeerWithXor[C Contact] []PeerWithXor[C]

func (l ListPeerWithXor[C]) Len() int {
	return len(l)
}
func (l ListPeerWithXor[C]) Slice(s, e int) ListPeerWithXor[C] {
	return l[s:e]
}

func (l ListPeerWithXor[C]) Less(i, j int) bool {
	return bytes.Compare(l[i].peer.UUID(), l[j].peer.UUID()) < 0
}

func (l ListPeerWithXor[C]) Swap(i, j int) {
	item := l[i]
	l[i] = l[j]
	l[j] = item
}

func (l *List[C]) Xor(peer C) ListPeerWithXor[C] {
	target := peer.Hash()
	result := make(ListPeerWithXor[C], 0, l.set.Len())
	for iter := l.set.Iter(); iter.HasNext(); iter.Next() {
		result = append(result, PeerWithXor[C]{peer: iter.Item().wrapped, xor: xor(target, iter.Item().wrapped.Hash())})
	}
	return result
}
