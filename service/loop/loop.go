package loop

type Element interface {
	Next() Element
}

func Loop(e Element) bool {
	var h = e
	var b bool
	for e != nil {
		e = e.Next()
		if h == e {
			return true
		}
		b = !b
		if b {
			h = h.Next()
		}
	}
	return false
}
