package storage

import "sync"

type MapIntersection struct {
	sync.RWMutex
	ids map[int64]map[int32]struct{}
}

func (s *MapIntersection) Insert(id int64, ipv4 int32) bool {
	s.Lock()
	defer s.Unlock()
	if s.ids == nil {
		s.ids = map[int64]map[int32]struct{}{id: {ipv4: struct{}{}}}
	} else {
		ips, ok := s.ids[id]
		if ok {
			_, ok := ips[ipv4]
			if ok {
				return false
			}
			ips[ipv4] = struct{}{}
		} else {
			s.ids[id] = map[int32]struct{}{ipv4: {}}
		}
	}
	return true
}

func (s *MapIntersection) Intersected(id1, id2 int64) bool {
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
	for ipv4 := range ips1 {
		_, ok := ips2[ipv4]
		if ok {
			n++
			if n == 2 {
				return true
			}
		}
	}
	return false
}
