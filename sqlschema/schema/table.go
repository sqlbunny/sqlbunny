package schema

// Table represents a database table.
type Table struct {
	Columns map[string]*Column `json:"columns"`

	PrimaryKey  *PrimaryKey            `json:"primary_key"`
	Indexes     map[string]*Index      `json:"indexes"`
	Uniques     map[string]*Unique     `json:"uniques"`
	ForeignKeys map[string]*ForeignKey `json:"foreign_keys"`
	Statistics  map[string]*Statistics `json:"statistics"`

	// StorageParams are the table's Postgres storage parameters
	// (WITH (name = value)). Values are SQL literals.
	StorageParams map[string]string `json:"storage_params"`

	// Unlogged is true for a Postgres UNLOGGED table. Writes to it skip the
	// WAL, and a crash empties it.
	Unlogged bool `json:"unlogged"`
}

func NewTable() *Table {
	return &Table{
		Columns:     make(map[string]*Column),
		Indexes:     make(map[string]*Index),
		Uniques:     make(map[string]*Unique),
		ForeignKeys: make(map[string]*ForeignKey),

		Statistics:    make(map[string]*Statistics),
		StorageParams: make(map[string]string),
	}
}
