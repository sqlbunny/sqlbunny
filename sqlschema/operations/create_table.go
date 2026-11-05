package operations

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

type Column struct {
	Name     string
	Type     string
	Default  string
	Nullable bool
}

type CreateTable struct {
	SchemaName string
	TableName  string
	Columns    []Column
	// Unlogged creates the table with CREATE UNLOGGED TABLE. The zero value
	// creates a normal (logged) table.
	Unlogged bool
	// StorageParams creates the table WITH (name = value, ...). Values are
	// SQL literals.
	StorageParams map[string]string
}

func (o CreateTable) GetSQL(transaction bool) string {
	var x []string
	for _, c := range o.Columns {
		var n string
		if !c.Nullable {
			n = " NOT NULL"
		}
		var d string
		if c.Default != "" {
			d = " DEFAULT " + c.Default
		}
		x = append(x, fmt.Sprintf("    \"%s\" %s%s%s", c.Name, c.Type, n, d))
	}
	kind := "TABLE"
	if o.Unlogged {
		kind = "UNLOGGED TABLE"
	}
	with := ""
	if len(o.StorageParams) != 0 {
		var params []string
		for _, k := range slices.Sorted(maps.Keys(o.StorageParams)) {
			params = append(params, fmt.Sprintf("%s = %s", k, o.StorageParams[k]))
		}
		with = " WITH (" + strings.Join(params, ", ") + ")"
	}
	return fmt.Sprintf("CREATE %s %s (\n%s\n)%s", kind, sqlName(o.SchemaName, o.TableName), strings.Join(x, ",\n"), with)
}

func (o CreateTable) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	if _, ok := s.Tables[o.TableName]; ok {
		return fmt.Errorf("table already exists: %s", o.TableName)
	}

	t := schema.NewTable()
	t.Unlogged = o.Unlogged
	for k, v := range o.StorageParams {
		t.StorageParams[k] = v
	}
	for _, c := range o.Columns {
		t.Columns[c.Name] = &schema.Column{
			Nullable: c.Nullable,
			Type:     c.Type,
			Default:  c.Default,
		}
	}
	s.Tables[o.TableName] = t

	return nil
}
