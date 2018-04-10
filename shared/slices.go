package shared

import (
	"errors"
	"fmt"
)

type StringSlice []string

func (s *StringSlice) Append(values ...string) {
	*s = append(*s, values...)
}

func (s *StringSlice) Insert(idx int, value string) error {
	if !s.isIndexValid(idx) {
		return errors.New("insertion failed. 'idx' bounds out of range")
	}
	*s = append(*s, "")
	copy((*s)[idx+1:], (*s)[idx:])
	(*s)[idx] = value
	return nil
}

func (s *StringSlice) Index(value string) (int, error) {
	for i, v := range *s {
		if v == value {
			return i, nil
		}
	}
	return 0, fmt.Errorf("'%s' not in slice", value)
}

func (s *StringSlice) DeleteByIndex(idx int) error {
	if !s.isIndexValid(idx) {
		return fmt.Errorf("deleting index failed. %d bounds out of range", idx)
	}
	*s = append((*s)[:idx], (*s)[idx+1:]...)
	return nil
}

func (s *StringSlice) DeleteByID(value string) error {
	i, err := s.Index(value)
	if err != nil {
		return err
	}
	return s.DeleteByIndex(i)
}

func (s *StringSlice) isIndexValid(idx int) bool {
	return idx >= 0 && idx < len(*s)
}
