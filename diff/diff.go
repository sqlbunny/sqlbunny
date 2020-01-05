package diff

import (
	"reflect"

	"github.com/sqlbunny/sqlschema/operations"
	"github.com/sqlbunny/sqlschema/schema"
)

func Diff(s1, s2 *schema.Schema) []operations.Operation {
	var ops []operations.Operation
	ops = diffDropForeignKeys(ops, s1, s2)
	ops = diffDropConstraints(ops, s1, s2)
	ops = diffDropIndexes(ops, s1, s2)
	ops = diffDropTables(ops, s1, s2)
	ops = diffAlterTables(ops, s1, s2)
	ops = diffCreateTables(ops, s1, s2)
	ops = diffCreateIndexes(ops, s1, s2)
	ops = diffCreateConstraints(ops, s1, s2)
	ops = diffCreateForeignKeys(ops, s1, s2)

	return ops
}

func hasPrimaryKey(s *schema.Schema, tableName string, k *schema.PrimaryKey) bool {
	t, ok := s.Tables[tableName]
	if !ok {
		return false
	}
	return reflect.DeepEqual(k, t.PrimaryKey)
}

func hasIndex(s *schema.Schema, tableName string, name string, k *schema.Index) bool {
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

func hasUnique(s *schema.Schema, tableName string, name string, k *schema.Unique) bool {
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

func hasForeignKey(s *schema.Schema, tableName string, name string, k *schema.ForeignKey) bool {
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

func diffDropForeignKeys(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for tableName, t1 := range s1.Tables {
		var subops []operations.AlterTableSuboperation

		for name, k := range t1.ForeignKeys {
			if !hasForeignKey(s2, tableName, name, k) {
				subops = append(subops, &operations.AlterTableDropForeignKey{Name: name})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, operations.AlterTable{
				TableName: tableName,
				Ops:       subops,
			})
		}
	}
	return ops
}

func diffDropConstraints(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for tableName, t1 := range s1.Tables {
		var subops []operations.AlterTableSuboperation

		if t1.PrimaryKey != nil && !hasPrimaryKey(s2, tableName, t1.PrimaryKey) {
			subops = append(subops, operations.AlterTableDropPrimaryKey{})
		}

		for name, k := range t1.Uniques {
			if !hasUnique(s2, tableName, name, k) {
				subops = append(subops, operations.AlterTableDropUnique{Name: name})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, operations.AlterTable{
				TableName: tableName,
				Ops:       subops,
			})
		}
	}
	return ops
}

func diffDropIndexes(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for tableName, t1 := range s1.Tables {
		for name, k := range t1.Indexes {
			if !hasIndex(s2, tableName, name, k) {
				ops = append(ops, operations.DropIndex{
					TableName: tableName,
					IndexName: name,
				})
			}
		}
	}
	return ops
}

func diffDropTables(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for name := range s1.Tables {
		if _, ok := s2.Tables[name]; !ok {
			ops = append(ops, operations.DropTable{
				TableName: name,
			})
		}
	}
	return ops
}

func diffAlterTables(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for tableName, t1 := range s1.Tables {
		if t2, ok := s2.Tables[tableName]; ok {
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
					TableName: tableName,
					Ops:       subops,
				})
			}

		}
	}
	return ops
}

func diffCreateTables(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for tableName, t2 := range s2.Tables {
		if _, ok := s1.Tables[tableName]; !ok {
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
				TableName: tableName,
				Columns:   cols,
			})
		}
	}
	return ops
}

func diffCreateIndexes(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for tableName, t2 := range s2.Tables {
		for name, i2 := range t2.Indexes {
			if !hasIndex(s1, tableName, name, i2) {
				ops = append(ops, operations.CreateIndex{
					TableName: tableName,
					IndexName: name,
					Columns:   i2.Columns,
				})
			}
		}
	}
	return ops
}

func diffCreateConstraints(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for tableName, t2 := range s2.Tables {
		var subops []operations.AlterTableSuboperation

		if t2.PrimaryKey != nil && !hasPrimaryKey(s1, tableName, t2.PrimaryKey) {
			subops = append(subops, operations.AlterTableCreatePrimaryKey{Columns: t2.PrimaryKey.Columns})
		}

		for name, i2 := range t2.Uniques {
			if !hasUnique(s1, tableName, name, i2) {
				subops = append(subops, operations.AlterTableCreateUnique{
					Name:    name,
					Columns: i2.Columns,
				})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, operations.AlterTable{
				TableName: tableName,
				Ops:       subops,
			})
		}
	}
	return ops
}

func diffCreateForeignKeys(ops []operations.Operation, s1, s2 *schema.Schema) []operations.Operation {
	for tableName, t2 := range s2.Tables {
		var subops []operations.AlterTableSuboperation

		for name, i2 := range t2.ForeignKeys {
			if !hasForeignKey(s1, tableName, name, i2) {
				subops = append(subops, operations.AlterTableCreateForeignKey{
					Name:           name,
					Columns:        i2.LocalColumns,
					ForeignTable:   i2.ForeignTable,
					ForeignColumns: i2.ForeignColumns,
				})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, operations.AlterTable{
				TableName: tableName,
				Ops:       subops,
			})
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
