package xorcontacts

import (
	"bytes"
	"errors"
	"github.com/itsabgr/go-set"
	"sort"
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
func (l *List[C]) GC() error {
	var todo []*contact[C]
	for iter := l.set.Iter(); iter.HasNext(); iter.Next() {
		if !iter.Item().Expired() {
			continue
		}
		todo = append(todo, iter.Item())
	}
	for _, item := range todo {
		l.set.Remove(item)
	}
	return nil
}

func (l *List[C]) Has(peer C) bool {
	return l.set.Has(&contact[C]{wrapped: peer})
}
func (l *List[C]) Len() int {
	return l.set.Len()
}
func (l *List[C]) List() *set.Set[*contact[C]] {
	return &l.set
}

type peerWithXor[C Contact] struct {
	peer C
	xor  []byte
}

func (p *peerWithXor[C]) Peer() C {
	return p.peer
}

func (p *peerWithXor[C]) Xor() []byte {
	return p.xor
}

type listPeerWithXor[C Contact] []peerWithXor[C]

func (l listPeerWithXor[C]) Len() int {
	return len(l)
}
func (l listPeerWithXor[C]) Slice(s, e int) listPeerWithXor[C] {
	return l[s:e]
}

func (l listPeerWithXor[C]) Less(i, j int) bool {
	return bytes.Compare(l[i].xor, l[j].xor) < 0
}

func (l listPeerWithXor[C]) Swap(i, j int) {
	item := l[i]
	l[i] = l[j]
	l[j] = item
}

func (l *List[C]) xor(peer C) listPeerWithXor[C] {
	target := hash(peer.ID())
	result := make(listPeerWithXor[C], 0, l.set.Len())
	for iter := l.set.Iter(); iter.HasNext(); iter.Next() {
		if iter.Item().Expired() {
			continue
		}
		result = append(result, peerWithXor[C]{peer: iter.Item().wrapped, xor: xor(target, hash(iter.Item().wrapped.ID()))})
	}
	return result
}
func (l *List[C]) Xor(peer C) []C {
	list := l.xor(peer)
	sort.Sort(list)
	result := make([]C, 0, len(list))
	for _, item := range list {
		result = append(result, item.peer)
	}
	return result
}
