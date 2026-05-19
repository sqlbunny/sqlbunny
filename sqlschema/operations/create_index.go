package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

type CreateIndex struct {
	SchemaName string
	TableName  string
	IndexName  string
	Columns    []string
	Method     string // Index method. If empty, default is btree.
	Where      string // Index where clause, for partial indexes. If empty, no where clause is in effect.
	Unique     bool   // If true, emits CREATE UNIQUE INDEX.
}

func (o CreateIndex) GetSQL() string {
	var buf bytes.Buffer
	kind := "INDEX"
	if o.Unique {
		kind = "UNIQUE INDEX"
	}
	fmt.Fprintf(&buf, "CREATE %s CONCURRENTLY \"%s\" ON %s", kind, o.IndexName, sqlName(o.SchemaName, o.TableName))
	if o.Method != "" {
		fmt.Fprintf(&buf, " USING %s", o.Method)
	}
	fmt.Fprintf(&buf, " (%s)", columnList(o.Columns))
	if o.Where != "" {
		fmt.Fprintf(&buf, " WHERE %s", o.Where)
	}
	return buf.String()
}

func (o CreateIndex) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	if _, ok := t.Indexes[o.IndexName]; ok {
		return fmt.Errorf("index already exists: %s ", o.IndexName)
	}
	t.Indexes[o.IndexName] = &schema.Index{
		Columns: o.Columns,
		Method:  o.Method,
		Where:   o.Where,
		Unique:  o.Unique,
	}
	return nil
}
