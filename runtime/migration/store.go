package migration

import (
	"fmt"
	"slices"

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
		if m == nil {
			return errors.Errorf("Migration '%s' is applied in the database, but is not in the migration store", mn)
		}
		for _, dn := range m.Dependencies {
			if _, ok := applied[dn]; !ok {
				return errors.Errorf("Migration '%s' is applied, but its dependency '%s' is not", mn, dn)
			}
		}
	}
	return nil
}

func (s *Store) validateMigrated(applied map[string]struct{}) error {
	for mn, _ := range s.Migrations {
		_, ok := applied[mn]
		if !ok {
			return errors.Errorf("Migration '%s' is in the migration store, but has not been applied to the database", mn)
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
	// Sorted, so a merge migration's dependency list doesn't depend on
	// map iteration order.
	slices.Sort(res)
	return res
}

func (s *Store) RunMigration(target string, applied map[string]struct{}, fn func(*Migration) error) error {
	if _, ok := applied[target]; ok {
		return nil
	}

	reverse := s.calcReverseDeps()

	// ready is kept in sorted batches rather than as a set, so migrations
	// are applied in the same order every run.
	var ready []string
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
			ready = append(ready, m.Name)
		} else {
			blocked[m.Name] = count
		}
	}
	slices.Sort(ready)

	for len(ready) != 0 {
		mn := ready[0]
		ready = ready[1:]

		m := s.Migrations[mn]
		if err := fn(m); err != nil {
			return err
		}

		// reverse[mn] comes out of a map, so sort what it unblocks before
		// queueing it.
		var unblocked []string
		for _, dn := range reverse[mn] {
			count, ok := blocked[dn]
			if !ok {
				continue
			}
			count--
			if count == 0 {
				unblocked = append(unblocked, dn)
				delete(blocked, dn)
			} else {
				blocked[dn] = count
			}
		}
		slices.Sort(unblocked)
		ready = append(ready, unblocked...)
	}

	if len(blocked) != 0 {
		panic("this should never happen")
	}

	return nil
}
