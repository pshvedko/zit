package storage

import (
	"slices"
	"sync"
)

type ArrayIntersection struct {
	sync.RWMutex
	ids map[int64][]int32
}

func (s *ArrayIntersection) Insert(id int64, ipv4 int32) bool {
	s.Lock()
	defer s.Unlock()
	if s.ids == nil {
		s.ids = map[int64][]int32{id: {ipv4}}
	} else {
		ips, ok := s.ids[id]
		if ok {
			n, ok := slices.BinarySearch(ips, ipv4)
			if ok {
				return false
			}
			s.ids[id] = slices.Insert(ips, n, ipv4)
		} else {
			s.ids[id] = []int32{ipv4}
		}
	}
	return true
}

func (s *ArrayIntersection) Intersected(id1, id2 int64) bool {
	if id1 == id2 {
		return true
	}
	s.RLock()
	defer s.RUnlock()
	ips1, ips2 := s.ids[id1], s.ids[id2]
	if len(ips1) > len(ips2) {
		ips1, ips2 = ips2, ips1
	}
	var n int
	for _, ipv4 := range ips1 {
		if len(ips2) == 0 {
			break
		}
		i, ok := slices.BinarySearch(ips2, ipv4)
		if ok {
			i++
			n++
			if n == 2 {
				return true
			}
		}
		ips2 = ips2[i:]
	}
	return false
}
