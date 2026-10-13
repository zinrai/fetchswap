package main

import "errors"

// Refuse, not let the last one win: a second -dest would silently replace
// another file than the first one names.
type onceString struct {
	value string
	set   bool
}

func (s *onceString) String() string { return s.value }

func (s *onceString) Set(v string) error {
	if s.set {
		return errors.New("given more than once")
	}
	s.value, s.set = v, true
	return nil
}
