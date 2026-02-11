package operations

import (
	"fmt"

	"github.com/sqlbunny/sqlbunny/sqlschema/schema"
)

type SetTableSchema struct {
	SchemaName    string
	TableName     string
	NewSchemaName string
}

func (o SetTableSchema) GetSQL() string {
	return fmt.Sprintf("ALTER TABLE %s SET SCHEMA \"%s\"", sqlName(o.SchemaName, o.TableName), o.NewSchemaName)
}

func (o SetTableSchema) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	s2, ok := d.Schemas[o.NewSchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.NewSchemaName)
	}
	t, ok := s.Tables[o.TableName]
	if !ok {
		return fmt.Errorf("no such table: %s", o.TableName)
	}
	if _, ok := s2.Tables[o.TableName]; ok {
		return fmt.Errorf("table with same name already exists in destination schema: %s", o.TableName)
	}

	delete(s.Tables, o.TableName)
	s2.Tables[o.TableName] = t

	for _, t2 := range s.Tables {
		for _, fk := range t2.ForeignKeys {
			if fk.ForeignSchema == o.SchemaName && fk.ForeignTable == o.TableName {
				fk.ForeignSchema = o.NewSchemaName
			}
		}
	}
	return nil
}
