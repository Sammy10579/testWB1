package user

import (
	"sync"
)

type Storage struct {
	m sync.Map
}

func (s *Storage) Set(ug UserGrade) {
	s.m.Store(ug.UserId, ug)
}

func (s *Storage) Get(UserId string) (UserGrade, error) {
	v, ok := s.m.Load(UserId)
	if ok {
		return v.(UserGrade), nil
	}
	return UserGrade{}, ErrNotFound
}
