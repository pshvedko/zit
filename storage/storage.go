package storage

type Storage map[int64]map[int32]struct{}

func (s Storage) Append(id int64, ip int32) error {
	ips, ok := s[id]
	if !ok {
		s[id] = map[int32]struct{}{ip: {}}
	} else {
		ips[ip] = struct{}{}
	}
	return nil
}

func (s Storage) Intersected(id1 int64, id2 int64) bool {
	ips1, ips2 := s[id1], s[id2]
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
