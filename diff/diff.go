package diff

import (
	"reflect"

	"github.com/sqlbunny/sqlschema/operations"
	"github.com/sqlbunny/sqlschema/schema"
)

func Diff(d1, d2 *schema.Database) []operations.Operation {
	var ops []operations.Operation
	ops = diffDropForeignKeys(ops, d1, d2)
	ops = diffDropConstraints(ops, d1, d2)
	ops = diffDropIndexes(ops, d1, d2)
	ops = diffDropTables(ops, d1, d2)
	ops = diffDropSchemas(ops, d1, d2)
	ops = diffAlterTables(ops, d1, d2)
	ops = diffCreateSchemas(ops, d1, d2)
	ops = diffCreateTables(ops, d1, d2)
	ops = diffCreateIndexes(ops, d1, d2)
	ops = diffCreateConstraints(ops, d1, d2)
	ops = diffCreateForeignKeys(ops, d1, d2)

	return ops
}

func getTable(d *schema.Database, schemaName, tableName string) *schema.Table {
	s, ok := d.Schemas[schemaName]
	if !ok {
		return nil
	}
	return s.Tables[tableName]
}

func hasPrimaryKey(d *schema.Database, schemaName, tableName string, k *schema.PrimaryKey) bool {
	s, ok := d.Schemas[schemaName]
	if !ok {
		return false
	}
	t, ok := s.Tables[tableName]
	if !ok {
		return false
	}
	return reflect.DeepEqual(k, t.PrimaryKey)
}

func hasIndex(d *schema.Database, schemaName, tableName string, name string, k *schema.Index) bool {
	s, ok := d.Schemas[schemaName]
	if !ok {
		return false
	}
	t, ok := s.Tables[tableName]
	if !ok {
		return false
	}
	k2, ok := t.Indexes[name]
	if !ok {
		return false
	}
	return reflect.DeepEqual(k, k2)
}

func hasUnique(d *schema.Database, schemaName, tableName string, name string, k *schema.Unique) bool {
	s, ok := d.Schemas[schemaName]
	if !ok {
		return false
	}
	t, ok := s.Tables[tableName]
	if !ok {
		return false
	}
	k2, ok := t.Uniques[name]
	if !ok {
		return false
	}
	return reflect.DeepEqual(k, k2)
}

func hasForeignKey(d *schema.Database, schemaName, tableName string, name string, k *schema.ForeignKey) bool {
	s, ok := d.Schemas[schemaName]
	if !ok {
		return false
	}
	t, ok := s.Tables[tableName]
	if !ok {
		return false
	}
	k2, ok := t.ForeignKeys[name]
	if !ok {
		return false
	}
	return reflect.DeepEqual(k, k2)
}

func diffDropForeignKeys(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s1 := range d1.Schemas {
		for tableName, t1 := range s1.Tables {
			var subops []operations.AlterTableSuboperation

			for name, k := range t1.ForeignKeys {
				if !hasForeignKey(d2, schemaName, tableName, name, k) {
					subops = append(subops, &operations.AlterTableDropForeignKey{Name: name})
				}
			}

			if len(subops) != 0 {
				ops = append(ops, operations.AlterTable{
					SchemaName: schemaName,
					TableName:  tableName,
					Ops:        subops,
				})
			}
		}
	}
	return ops
}

func diffDropConstraints(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s1 := range d1.Schemas {
		for tableName, t1 := range s1.Tables {
			var subops []operations.AlterTableSuboperation

			if t1.PrimaryKey != nil && !hasPrimaryKey(d2, schemaName, tableName, t1.PrimaryKey) {
				subops = append(subops, operations.AlterTableDropPrimaryKey{})
			}

			for name, k := range t1.Uniques {
				if !hasUnique(d2, schemaName, tableName, name, k) {
					subops = append(subops, operations.AlterTableDropUnique{Name: name})
				}
			}

			if len(subops) != 0 {
				ops = append(ops, operations.AlterTable{
					SchemaName: schemaName,
					TableName:  tableName,
					Ops:        subops,
				})
			}
		}
	}
	return ops
}

func diffDropIndexes(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s1 := range d1.Schemas {
		for tableName, t1 := range s1.Tables {
			for name, k := range t1.Indexes {
				if !hasIndex(d2, schemaName, tableName, name, k) {
					ops = append(ops, operations.DropIndex{
						SchemaName: schemaName,
						TableName:  tableName,
						IndexName:  name,
					})
				}
			}
		}
	}
	return ops
}

func diffDropTables(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s1 := range d1.Schemas {
		for tableName := range s1.Tables {
			if getTable(d2, schemaName, tableName) == nil {
				ops = append(ops, operations.DropTable{
					SchemaName: schemaName,
					TableName:  tableName,
				})
			}
		}
	}
	return ops
}

func diffDropSchemas(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName := range d1.Schemas {
		if d2.Schemas[schemaName] == nil {
			ops = append(ops, operations.DropSchema{
				SchemaName: schemaName,
			})
		}
	}
	return ops
}

func diffAlterTables(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s1 := range d1.Schemas {
		for tableName, t1 := range s1.Tables {
			if t2 := getTable(d2, schemaName, tableName); t2 != nil {
				var subops []operations.AlterTableSuboperation
				for name, c1 := range t1.Columns {
					c2, ok := t2.Columns[name]
					if ok {
						subops = diffColumn(subops, name, c1, c2)
					} else {
						subops = append(subops, operations.AlterTableDropColumn{Name: name})
					}
				}
				for name, c2 := range t2.Columns {
					_, ok := t1.Columns[name]
					if !ok {
						subops = append(subops, operations.AlterTableAddColumn{
							Name:     name,
							Type:     c2.Type,
							Default:  c2.Default,
							Nullable: c2.Nullable,
						})
					}
				}

				if len(subops) != 0 {
					ops = append(ops, operations.AlterTable{
						SchemaName: schemaName,
						TableName:  tableName,
						Ops:        subops,
					})
				}
			}
		}
	}
	return ops
}

func diffCreateSchemas(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName := range d2.Schemas {
		if d1.Schemas[schemaName] == nil {
			ops = append(ops, operations.CreateSchema{
				SchemaName: schemaName,
			})
		}
	}
	return ops
}

func diffCreateTables(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s2 := range d2.Schemas {
		for tableName, t2 := range s2.Tables {
			if t1 := getTable(d1, schemaName, tableName); t1 == nil {
				var cols []operations.Column
				for name, c := range t2.Columns {
					cols = append(cols, operations.Column{
						Name:     name,
						Type:     c.Type,
						Default:  c.Default,
						Nullable: c.Nullable,
					})
				}

				ops = append(ops, operations.CreateTable{
					SchemaName: schemaName,
					TableName:  tableName,
					Columns:    cols,
				})
			}
		}
	}
	return ops
}

func diffCreateIndexes(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s2 := range d2.Schemas {
		for tableName, t2 := range s2.Tables {
			for name, i2 := range t2.Indexes {
				if !hasIndex(d1, schemaName, tableName, name, i2) {
					ops = append(ops, operations.CreateIndex{
						SchemaName: schemaName,
						TableName:  tableName,
						IndexName:  name,
						Columns:    i2.Columns,
					})
				}
			}
		}
	}
	return ops
}

func diffCreateConstraints(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s2 := range d2.Schemas {
		for tableName, t2 := range s2.Tables {
			var subops []operations.AlterTableSuboperation

			if t2.PrimaryKey != nil && !hasPrimaryKey(d1, schemaName, tableName, t2.PrimaryKey) {
				subops = append(subops, operations.AlterTableCreatePrimaryKey{Columns: t2.PrimaryKey.Columns})
			}

			for name, i2 := range t2.Uniques {
				if !hasUnique(d1, schemaName, tableName, name, i2) {
					subops = append(subops, operations.AlterTableCreateUnique{
						Name:    name,
						Columns: i2.Columns,
					})
				}
			}

			if len(subops) != 0 {
				ops = append(ops, operations.AlterTable{
					SchemaName: schemaName,
					TableName:  tableName,
					Ops:        subops,
				})
			}
		}
	}
	return ops
}

func diffCreateForeignKeys(ops []operations.Operation, d1, d2 *schema.Database) []operations.Operation {
	for schemaName, s2 := range d2.Schemas {
		for tableName, t2 := range s2.Tables {
			var subops []operations.AlterTableSuboperation

			for name, i2 := range t2.ForeignKeys {
				if !hasForeignKey(d1, schemaName, tableName, name, i2) {
					subops = append(subops, operations.AlterTableCreateForeignKey{
						Name:           name,
						Columns:        i2.LocalColumns,
						ForeignSchema:  i2.ForeignSchema,
						ForeignTable:   i2.ForeignTable,
						ForeignColumns: i2.ForeignColumns,
					})
				}
			}

			if len(subops) != 0 {
				ops = append(ops, operations.AlterTable{
					SchemaName: schemaName,
					TableName:  tableName,
					Ops:        subops,
				})
			}
		}
	}
	return ops
}

func diffColumn(ops []operations.AlterTableSuboperation, name string, c1, c2 *schema.Column) []operations.AlterTableSuboperation {
	if c1.Nullable && !c2.Nullable {
		ops = append(ops, operations.AlterTableSetNotNull{Name: name})
	}
	if !c1.Nullable && c2.Nullable {
		ops = append(ops, operations.AlterTableSetNull{Name: name})
	}
	if c1.Default != c2.Default {
		if c2.Default == "" {
			ops = append(ops, operations.AlterTableDropDefault{
				Name: name,
			})
		} else {
			ops = append(ops, operations.AlterTableSetDefault{
				Name:    name,
				Default: c2.Default,
			})
		}
	}
	if c1.Type != c2.Type {
		ops = append(ops, operations.AlterTableSetType{
			Name: name,
			Type: c2.Type,
		})
	}
	return ops
}
