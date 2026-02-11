package schema

// Table represents a database table.
type Table struct {
	Columns map[string]*Column `json:"columns"`

	PrimaryKey  *PrimaryKey            `json:"primary_key"`
	Indexes     map[string]*Index      `json:"indexes"`
	Uniques     map[string]*Unique     `json:"uniques"`
	ForeignKeys map[string]*ForeignKey `json:"foreign_keys"`
}

func NewTable() *Table {
	return &Table{
		Columns:     make(map[string]*Column),
		Indexes:     make(map[string]*Index),
		Uniques:     make(map[string]*Unique),
		ForeignKeys: make(map[string]*ForeignKey),
	}
}
