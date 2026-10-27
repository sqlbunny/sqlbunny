package operations

import (
	"fmt"
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
	return fmt.Sprintf("CREATE %s %s (\n%s\n)", kind, sqlName(o.SchemaName, o.TableName), strings.Join(x, ",\n"))
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
