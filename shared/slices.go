package shared

type StringSlice []string

func (s *StringSlice) Append(values ...string) {
	*s = append(*s, values...)
}

func (s *StringSlice) Insert(idx int, value string) bool {
	if !s.isIndexValid(idx) {
		return false
	}
	*s = append(*s, "")
	copy((*s)[idx+1:], (*s)[idx:])
	(*s)[idx] = value
	return true
}

func (s *StringSlice) ExpandAfterID(id string, values ...string) bool {
	idx := s.Index(id)
	if !s.isIndexValid(idx) {
		return false
	}
	*s = append((*s)[:idx+1], append(values, (*s)[idx+1:]...)...)
	return true
}

func (s *StringSlice) ExpandBeforeID(id string, values ...string) bool {
	idx := s.Index(id)
	if !s.isIndexValid(idx) {
		return false
	}
	*s = append((*s)[:idx], append(values, (*s)[idx:]...)...)
	return true
}

func (s *StringSlice) InsertBeforeID(id, v string) bool {
	return s.Insert(s.Index(id), v)
}

func (s *StringSlice) InsertAfterID(id, v string) bool {
	if s.Include(id) {
		return s.Insert(s.Index(id)+1, v)
	}
	return false

}

func (s StringSlice) Index(value string) int {
	for i, v := range s {
		if v == value {
			return i
		}
	}
	return -1
}

func (s StringSlice) Include(value string) bool {
	return s.Index(value) >= 0

}

func (s *StringSlice) DeleteByIndex(idx int) bool {
	if !s.isIndexValid(idx) {
		return false
	}
	*s = append((*s)[:idx], (*s)[idx+1:]...)
	return true
}

func (s *StringSlice) DeleteByID(value string) bool {
	return s.DeleteByIndex(s.Index(value))
}

func (s *StringSlice) isIndexValid(idx int) bool {
	return idx >= 0 && idx < len(*s)
}
