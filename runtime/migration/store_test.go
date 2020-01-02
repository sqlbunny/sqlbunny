package migration

import "testing"

func equal(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func checkEqual(t *testing.T, what string, got, expected []string) {
	if !equal(got, expected) {
		t.Errorf("%s: expected %v, got %v", what, got, expected)
	}
}

func equalUnsorted(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	// create a map of string -> int
	diff := make(map[string]int, len(x))
	for _, sx := range x {
		// 0 value for int is 0, so just increment a counter for the string
		diff[sx]++
	}
	for _, sy := range y {
		// If the string _y is not in diff bail out early
		if _, ok := diff[sy]; !ok {
			return false
		}
		diff[sy]--
		if diff[sy] == 0 {
			delete(diff, sy)
		}
	}
	return len(diff) == 0
}

func checkEqualUnsorted(t *testing.T, what string, got, expected []string) {
	if !equalUnsorted(got, expected) {
		t.Errorf("%s: expected %v, got %v", what, got, expected)
	}
}

func TestStoreGetHeads(t *testing.T) {
	s := Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name: "a",
			},
			"b": &Migration{
				Name: "b",
			},
			"c": &Migration{
				Name: "c",
			},
		},
	}
	checkEqualUnsorted(t, "no deps", s.FindHeads(), []string{"a", "b", "c"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name:         "a",
				Dependencies: []string{"b"},
			},
			"b": &Migration{
				Name:         "b",
				Dependencies: []string{"c"},
			},
			"c": &Migration{
				Name: "c",
			},
		},
	}
	checkEqualUnsorted(t, "chain", s.FindHeads(), []string{"a"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name:         "a",
				Dependencies: []string{"b", "c"},
			},
			"b": &Migration{
				Name: "b",
			},
			"c": &Migration{
				Name: "c",
			},
		},
	}
	checkEqualUnsorted(t, "tree", s.FindHeads(), []string{"a"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name: "a",
			},
			"b": &Migration{
				Name:         "b",
				Dependencies: []string{"a", "c"},
			},
			"c": &Migration{
				Name: "c",
			},
		},
	}
	checkEqualUnsorted(t, "tree 2", s.FindHeads(), []string{"b"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name:         "a",
				Dependencies: []string{"b", "c"},
			},
			"b": &Migration{
				Name:         "b",
				Dependencies: []string{"d", "e"},
			},
			"c": &Migration{
				Name: "c",
			},
			"d": &Migration{
				Name: "d",
			},
			"e": &Migration{
				Name:         "e",
				Dependencies: []string{"f"},
			},
			"f": &Migration{
				Name: "f",
			},
			"g": &Migration{
				Name: "g",
			},
		},
	}
	checkEqualUnsorted(t, "tree deep", s.FindHeads(), []string{"a", "g"})
}

func makeCallback() (*[]string, func(*Migration) error) {
	res := new([]string)

	return res, func(m *Migration) error {
		*res = append(*res, m.Name)
		return nil
	}
}

func TestRun(t *testing.T) {
	s := Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name: "a",
			},
			"b": &Migration{
				Name: "b",
			},
			"c": &Migration{
				Name: "c",
			},
		},
	}
	r, f := makeCallback()
	err := s.RunMigration("a", map[string]struct{}{}, f)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	checkEqual(t, "no deps", *r, []string{"a"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name:         "a",
				Dependencies: []string{"b"},
			},
			"b": &Migration{
				Name:         "b",
				Dependencies: []string{"c"},
			},
			"c": &Migration{
				Name: "c",
			},
		},
	}
	r, f = makeCallback()
	err = s.RunMigration("a", map[string]struct{}{}, f)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	checkEqual(t, "chain", *r, []string{"c", "b", "a"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name:         "a",
				Dependencies: []string{"b"},
			},
			"b": &Migration{
				Name:         "b",
				Dependencies: []string{"c"},
			},
			"c": &Migration{
				Name: "c",
			},
		},
	}
	r, f = makeCallback()
	err = s.RunMigration("b", map[string]struct{}{}, f)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	checkEqual(t, "half chain", *r, []string{"c", "b"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name:         "a",
				Dependencies: []string{"b", "c", "d"},
			},
			"b": &Migration{
				Name: "b",
			},
			"c": &Migration{
				Name: "c",
			},
			"d": &Migration{
				Name: "d",
			},
		},
	}
	r, f = makeCallback()
	err = s.RunMigration("a", map[string]struct{}{}, f)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	checkEqualUnsorted(t, "tree", (*r)[0:3], []string{"d", "c", "b"})
	checkEqualUnsorted(t, "tree", (*r)[3:], []string{"a"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name:         "a",
				Dependencies: []string{"b", "c", "d"},
			},
			"b": &Migration{
				Name: "b",
			},
			"c": &Migration{
				Name: "c",
			},
			"d": &Migration{
				Name: "d",
			},
		},
	}
	r, f = makeCallback()
	err = s.RunMigration("a", map[string]struct{}{
		"b": {},
		"c": {},
	}, f)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	checkEqualUnsorted(t, "tree2", (*r)[0:1], []string{"d"})
	checkEqualUnsorted(t, "tree2", (*r)[1:], []string{"a"})

	s = Store{
		Migrations: map[string]*Migration{
			"a": &Migration{
				Name:         "a",
				Dependencies: []string{"b", "c", "d"},
			},
			"b": &Migration{
				Name: "b",
			},
			"c": &Migration{
				Name: "c",
			},
			"d": &Migration{
				Name: "d",
			},
		},
	}
	r, f = makeCallback()
	err = s.RunMigration("a", map[string]struct{}{
		"a": {},
		"b": {},
		"c": {},
		"d": {},
	}, f)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	checkEqualUnsorted(t, "tree3", *r, []string{})
}
