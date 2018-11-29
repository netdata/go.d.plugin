package utils

type StringSlice []string

func (s *StringSlice) Append(values ...string) {
	*s = append(*s, values...)
}

func (s *StringSlice) Insert(idx int, values ...string) bool {
	if !s.isIndexValid(idx) {
		return false
	}
	s.insert(idx, values...)
	return true
}

func (s *StringSlice) InsertAfterID(id string, values ...string) bool {
	idx := s.Index(id)
	if !s.isIndexValid(idx) {
		return false
	}
	s.insert(idx+1, values...)
	return true
}

func (s *StringSlice) InsertBeforeID(id string, values ...string) bool {
	return s.Insert(s.Index(id), values...)
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

func (s *StringSlice) insert(idx int, values ...string) {
	*s = append((*s)[:idx], append(values, (*s)[idx:]...)...)
}
