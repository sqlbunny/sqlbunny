package migration

import (
	"reflect"

	"github.com/sqlbunny/sqlbunny/runtime/migration"
	"github.com/sqlbunny/sqlbunny/schema"
)

func diff(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	ops = diffDropForeignKeys(ops, s1, s2)
	ops = diffDropConstraints(ops, s1, s2)
	ops = diffDropIndexes(ops, s1, s2)
	ops = diffDropModels(ops, s1, s2)
	ops = diffAlterModels(ops, s1, s2)
	ops = diffCreateModels(ops, s1, s2)
	ops = diffCreateIndexes(ops, s1, s2)
	ops = diffCreateConstraints(ops, s1, s2)
	ops = diffCreateForeignKeys(ops, s1, s2)

	return ops
}

func hasPrimaryKey(m *schema.Model, columns []string) bool {
	if m == nil {
		return false
	}
	return m.PrimaryKey != nil && reflect.DeepEqual(m.PrimaryKey.Columns, columns)
}

func hasUnique(m *schema.Model, indexName string) bool {
	if m == nil {
		return false
	}
	return m.FindUnique(indexName) != nil
}

func hasIndex(m *schema.Model, indexName string) bool {
	if m == nil {
		return false
	}
	return m.FindIndex(indexName) != nil
}

func hasForeignKey(m *schema.Model, indexName string) bool {
	if m == nil {
		return false
	}
	return m.FindForeignKey(indexName) != nil
}

func diffDropForeignKeys(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name, m1 := range s1.Models {
		var subops []migration.AlterTableSuboperation

		for _, i1 := range m1.ForeignKeys {
			if !hasForeignKey(s2.Models[name], i1.Name) {
				subops = append(subops, &migration.AlterTableDropForeignKey{Name: i1.Name})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, migration.AlterTableOperation{
				Name: m1.Name,
				Ops:  subops,
			})
		}
	}
	return ops
}

func diffDropConstraints(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name, m1 := range s1.Models {
		var subops []migration.AlterTableSuboperation

		if m1.PrimaryKey != nil && !hasPrimaryKey(s2.Models[name], m1.PrimaryKey.Columns) {
			subops = append(subops, migration.AlterTableDropPrimaryKey{})
		}

		for _, i1 := range m1.Uniques {
			if !hasUnique(s2.Models[name], i1.Name) {
				subops = append(subops, migration.AlterTableDropUnique{Name: i1.Name})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, migration.AlterTableOperation{
				Name: m1.Name,
				Ops:  subops,
			})
		}
	}
	return ops
}

func diffDropIndexes(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name, m1 := range s1.Models {
		for _, i1 := range m1.Indexes {
			if !hasIndex(s2.Models[name], i1.Name) {
				ops = append(ops, migration.DropIndexOperation{
					Name:      name,
					IndexName: i1.Name,
				})
			}
		}
	}
	return ops
}

func diffDropModels(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name := range s1.Models {
		if _, ok := s2.Models[name]; !ok {
			ops = append(ops, migration.DropTableOperation{
				Name: name,
			})
		}
	}
	return ops
}

func diffAlterModels(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name, m1 := range s1.Models {
		if m2, ok := s2.Models[name]; ok {
			var subops []migration.AlterTableSuboperation
			for _, c1 := range m1.Columns {
				c2 := m2.FindColumn(c1.Name)
				if c2 != nil {
					subops = diffColumn(subops, c1, c2)
				} else {
					subops = append(subops, migration.AlterTableDropColumn{Name: c1.Name})
				}
			}
			for _, c2 := range m2.Columns {
				c1 := m1.FindColumn(c2.Name)
				if c1 == nil {
					subops = append(subops, migration.AlterTableAddColumn{
						Name:     c2.Name,
						Type:     c2.SQLType,
						Default:  c2.SQLDefault,
						Nullable: c2.Nullable,
					})
				}
			}

			if len(subops) != 0 {
				ops = append(ops, migration.AlterTableOperation{
					Name: m1.Name,
					Ops:  subops,
				})
			}

		}
	}
	return ops
}

func diffCreateModels(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name, m2 := range s2.Models {
		if _, ok := s1.Models[name]; !ok {
			ops = append(ops, migration.CreateTableOperation{
				Name:    name,
				Columns: makeColumns(m2.Columns),
			})
		}
	}
	return ops
}

func diffCreateIndexes(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name, m2 := range s2.Models {
		for _, i2 := range m2.Indexes {
			if !hasIndex(s1.Models[name], i2.Name) {
				ops = append(ops, migration.CreateIndexOperation{
					Name:      name,
					IndexName: i2.Name,
					Columns:   i2.Columns,
				})
			}
		}
	}
	return ops
}

func diffCreateConstraints(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name, m2 := range s2.Models {
		var subops []migration.AlterTableSuboperation

		if m2.PrimaryKey != nil && !hasPrimaryKey(s1.Models[name], m2.PrimaryKey.Columns) {
			subops = append(subops, migration.AlterTableCreatePrimaryKey{Columns: m2.PrimaryKey.Columns})
		}

		for _, i2 := range m2.Uniques {
			if !hasUnique(s1.Models[name], i2.Name) {
				subops = append(subops, migration.AlterTableCreateUnique{
					Name:    i2.Name,
					Columns: i2.Columns,
				})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, migration.AlterTableOperation{
				Name: m2.Name,
				Ops:  subops,
			})
		}
	}
	return ops
}

func diffCreateForeignKeys(ops migration.OperationList, s1, s2 *schema.Schema) migration.OperationList {
	for name, m2 := range s2.Models {
		var subops []migration.AlterTableSuboperation

		for _, i2 := range m2.ForeignKeys {
			if !hasForeignKey(s1.Models[name], i2.Name) {
				subops = append(subops, migration.AlterTableCreateForeignKey{
					Name:           i2.Name,
					Columns:        i2.Columns,
					ForeignTable:   i2.ForeignModel,
					ForeignColumns: i2.ForeignColumns,
				})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, migration.AlterTableOperation{
				Name: m2.Name,
				Ops:  subops,
			})
		}
	}
	return ops
}

func makeColumns(m []*schema.Column) []migration.Column {
	var res []migration.Column
	for _, c := range m {
		res = append(res, migration.Column{
			Name:     c.Name,
			Type:     c.SQLType,
			Default:  c.SQLDefault,
			Nullable: c.Nullable,
		})
	}
	return res
}

func diffColumn(ops []migration.AlterTableSuboperation, c1, c2 *schema.Column) []migration.AlterTableSuboperation {
	if c1.Nullable && !c2.Nullable {
		ops = append(ops, migration.AlterTableSetNotNull{Name: c1.Name})
	}
	if !c1.Nullable && c2.Nullable {
		ops = append(ops, migration.AlterTableSetNull{Name: c1.Name})
	}
	if c1.SQLDefault != c2.SQLDefault {
		if c2.SQLDefault == "" {
			ops = append(ops, migration.AlterTableDropDefault{
				Name: c1.Name,
			})
		} else {
			ops = append(ops, migration.AlterTableSetDefault{
				Name:    c1.Name,
				Default: c2.SQLDefault,
			})
		}
	}
	if c1.SQLType != c2.SQLType {
		ops = append(ops, migration.AlterTableSetType{
			Name: c1.Name,
			Type: c2.SQLType,
		})
	}
	return ops
}
