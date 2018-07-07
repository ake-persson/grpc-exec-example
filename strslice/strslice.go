package strslice

import "strings"

type StrSlice []string

func (s *StrSlice) String() string {
	return strings.Join([]string(*s), ",")
}

func (s *StrSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *StrSlice) Len() int {
	return len(*s)
}

func (s *StrSlice) Index(i int) string {
	if i >= s.Len() {
		return ""
	}
	e := *s
	return e[i]
}

func (s *StrSlice) First() string {
	if s.Len() == 0 {
		return ""
	}
	return s.Index(0)
}

func (s *StrSlice) Append(v string) {
	*s = append(*s, v)
}
