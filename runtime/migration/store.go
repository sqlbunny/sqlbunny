package migration

import (
	"fmt"

	"github.com/sqlbunny/errors"
)

type Store struct {
	Migrations map[string]*Migration
}

func (s *Store) Register(m *Migration) {
	if s.Migrations == nil {
		s.Migrations = make(map[string]*Migration)
	}
	if _, ok := s.Migrations[m.Name]; ok {
		panic(fmt.Sprintf("Migration with name '%s' registered multiple times", m.Name))
	}
	s.Migrations[m.Name] = m
}

func (s *Store) Validate() error {
	// TODO: Check no migrations with name ""
	// TODO: Check no cycles
	// TODO: Check all migrations dependencies exist
	// TODO: Check that the map key and the object name match
	// TODO: Check that migration dependencies slices contains no duplicates
	// TODO: Check that migration dependencies slices doesn't contain the migration itself (special case of no cycles)
	return nil
}

func (s *Store) calcReverseDeps() map[string][]string {
	res := make(map[string][]string)
	for _, m := range s.Migrations {
		for _, d := range m.Dependencies {
			res[d] = append(res[d], m.Name)
		}
	}
	return res
}

func (s *Store) validateApplied(applied map[string]struct{}) error {
	for mn := range applied {
		m := s.Migrations[mn]
		for _, dn := range m.Dependencies {
			if _, ok := applied[dn]; !ok {
				return errors.Errorf("Migration '%s' is applied, but its dependency '%s' is not", mn, dn)
			}
		}
	}
	return nil
}

func (s *Store) FindHeads() []string {
	notHeads := make(map[string]struct{})
	for _, m := range s.Migrations {
		for _, dn := range m.Dependencies {
			notHeads[dn] = struct{}{}
		}
	}
	var res []string
	for _, m := range s.Migrations {
		if _, ok := notHeads[m.Name]; !ok {
			res = append(res, m.Name)
		}
	}
	return res
}

func (s *Store) RunMigration(target string, applied map[string]struct{}, fn func(*Migration) error) error {
	if _, ok := applied[target]; ok {
		return nil
	}

	reverse := s.calcReverseDeps()

	ready := make(map[string]struct{})
	blocked := make(map[string]int)

	visited := make(map[string]struct{})
	q := []string{target}
	for len(q) != 0 {
		m := s.Migrations[q[0]]
		q = q[1:]

		count := 0
		for _, dn := range m.Dependencies {
			if _, ok := applied[dn]; !ok {
				count++
				if _, ok2 := visited[dn]; !ok2 {
					visited[dn] = struct{}{}
					q = append(q, dn)
				}
			}
		}
		if count == 0 {
			ready[m.Name] = struct{}{}
		} else {
			blocked[m.Name] = count
		}
	}

	for len(ready) != 0 {
		for mn := range ready {
			delete(ready, mn)
			m := s.Migrations[mn]
			if err := fn(m); err != nil {
				return err
			}
			for _, dn := range reverse[mn] {
				count, ok := blocked[dn]
				if !ok {
					continue
				}
				count--
				if count == 0 {
					ready[dn] = struct{}{}
					delete(blocked, dn)
				} else {
					blocked[dn] = count
				}
			}
		}
	}

	if len(blocked) != 0 {
		panic("this should never happen")
	}

	return nil
}
