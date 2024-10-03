package operations

import (
	"bytes"
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type CreateIndex struct {
	SchemaName string
	TableName  string
	IndexName  string
	Columns    []string
	Method     string // Index method. If empty, default is btree.
	Where      string // Index where clause, for partial indexes. If empty, no where clause is in effect.
}

func (o CreateIndex) GetSQL() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "CREATE INDEX CONCURRENTLY \"%s\" ON %s", o.IndexName, sqlName(o.SchemaName, o.TableName))
	if o.Method != "" {
		fmt.Fprintf(&buf, " USING %s", o.Method)
	}
	fmt.Fprintf(&buf, " (%s)", columnList(o.Columns))
	if o.Where != "" {
		fmt.Fprintf(&buf, " WHERE %s", o.Where)
	}
	return buf.String()
}

func (o CreateIndex) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.CreateIndex {\n")
	fmt.Fprint(w, "SchemaName: "+esc(o.SchemaName)+",\n")
	fmt.Fprint(w, "TableName: "+esc(o.TableName)+",\n")
	fmt.Fprint(w, "IndexName: "+esc(o.IndexName)+",\n")
	fmt.Fprint(w, "Columns: []string{"+columnList(o.Columns)+"},\n")
	if o.Method != "" {
		fmt.Fprint(w, "Method: "+esc(o.Method)+",\n")
	}
	if o.Where != "" {
		fmt.Fprint(w, "Where: "+esc(o.Where)+",\n")
	}
	fmt.Fprint(w, "}")
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
	}
	return nil
}
