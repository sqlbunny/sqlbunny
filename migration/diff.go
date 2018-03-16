package migration

import (
	"reflect"

	"github.com/KernelPay/sqlboiler/schema"
)

func Diff(ops OperationList, s1, s2 *schema.Schema) OperationList {
	ops = DiffDropForeignKeys(ops, s1, s2)
	ops = DiffDropConstraints(ops, s1, s2)
	ops = DiffDropIndexes(ops, s1, s2)
	ops = DiffDropModels(ops, s1, s2)
	ops = DiffAlterModels(ops, s1, s2)
	ops = DiffCreateModels(ops, s1, s2)
	ops = DiffCreateIndexes(ops, s1, s2)
	ops = DiffCreateConstraints(ops, s1, s2)
	ops = DiffCreateForeignKeys(ops, s1, s2)

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

func DiffDropForeignKeys(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name, m1 := range s1.ModelsByName {
		var subops []AlterTableSuboperation

		for _, i1 := range m1.ForeignKeys {
			if !hasForeignKey(s2.ModelsByName[name], i1.Name) {
				subops = append(subops, &AlterTableDropForeignKey{Name: i1.Name})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, AlterTableOperation{
				Name: m1.Name,
				Ops:  subops,
			})
		}
	}
	return ops
}

func DiffDropConstraints(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name, m1 := range s1.ModelsByName {
		var subops []AlterTableSuboperation

		if m1.PrimaryKey != nil && !hasPrimaryKey(s2.ModelsByName[name], m1.PrimaryKey.Columns) {
			subops = append(subops, &AlterTableDropPrimaryKey{})
		}

		for _, i1 := range m1.Uniques {
			if !hasUnique(s2.ModelsByName[name], i1.Name) {
				subops = append(subops, &AlterTableDropUnique{Name: i1.Name})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, AlterTableOperation{
				Name: m1.Name,
				Ops:  subops,
			})
		}
	}
	return ops
}

func DiffDropIndexes(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name, m1 := range s1.ModelsByName {
		for _, i1 := range m1.Indexes {
			if !hasIndex(s2.ModelsByName[name], i1.Name) {
				ops = append(ops, DropIndexOperation{
					Name:      name,
					IndexName: i1.Name,
				})
			}
		}
	}
	return ops
}

func DiffDropModels(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name := range s1.ModelsByName {
		if _, ok := s2.ModelsByName[name]; !ok {
			ops = append(ops, DropTableOperation{
				Name: name,
			})
		}
	}
	return ops
}

func DiffAlterModels(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name, m1 := range s1.ModelsByName {
		if m2, ok := s2.ModelsByName[name]; ok {
			var subops []AlterTableSuboperation
			for _, c1 := range m1.Columns {
				c2 := m2.FindColumn(c1.Name)
				if c2 != nil {
					subops = DiffColumn(subops, c1, c2)
				} else {
					subops = append(subops, &AlterTableDropColumn{Name: c1.Name})
				}
			}
			for _, c2 := range m2.Columns {
				c1 := m1.FindColumn(c2.Name)
				if c1 == nil {
					subops = append(subops, &AlterTableAddColumn{
						Name:     c2.Name,
						Type:     c2.DBType,
						Nullable: c2.Nullable,
					})
				}
			}

			if len(subops) != 0 {
				ops = append(ops, AlterTableOperation{
					Name: m1.Name,
					Ops:  subops,
				})
			}

		}
	}
	return ops
}

func DiffCreateModels(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name, m2 := range s2.ModelsByName {
		if _, ok := s1.ModelsByName[name]; !ok {
			ops = append(ops, CreateTableOperation{
				Name:    name,
				Columns: makeColumns(m2.Columns),
			})
		}
	}
	return ops
}

func DiffCreateIndexes(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name, m2 := range s2.ModelsByName {
		for _, i2 := range m2.Indexes {
			if !hasIndex(s1.ModelsByName[name], i2.Name) {
				ops = append(ops, CreateIndexOperation{
					Name:      name,
					IndexName: i2.Name,
					Columns:   i2.Columns,
				})
			}
		}
	}
	return ops
}

func DiffCreateConstraints(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name, m2 := range s2.ModelsByName {
		var subops []AlterTableSuboperation

		if m2.PrimaryKey != nil && !hasPrimaryKey(s1.ModelsByName[name], m2.PrimaryKey.Columns) {
			subops = append(subops, &AlterTableCreatePrimaryKey{Columns: m2.PrimaryKey.Columns})
		}

		for _, i2 := range m2.Uniques {
			if !hasUnique(s1.ModelsByName[name], i2.Name) {
				subops = append(subops, &AlterTableCreateUnique{
					Name:    i2.Name,
					Columns: i2.Columns,
				})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, AlterTableOperation{
				Name: m2.Name,
				Ops:  subops,
			})
		}
	}
	return ops
}

func DiffCreateForeignKeys(ops OperationList, s1, s2 *schema.Schema) OperationList {
	for name, m2 := range s2.ModelsByName {
		var subops []AlterTableSuboperation

		for _, i2 := range m2.ForeignKeys {
			if !hasForeignKey(s1.ModelsByName[name], i2.Name) {
				subops = append(subops, &AlterTableCreateForeignKey{
					Name:           i2.Name,
					Columns:        []string{i2.Column},
					ForeignTable:   i2.ForeignModel,
					ForeignColumns: []string{i2.ForeignColumn},
				})
			}
		}

		if len(subops) != 0 {
			ops = append(ops, AlterTableOperation{
				Name: m2.Name,
				Ops:  subops,
			})
		}
	}
	return ops
}

func makeColumns(m []*schema.Column) []Column {
	var res []Column
	for _, c := range m {
		res = append(res, Column{
			Name:     c.Name,
			Type:     c.DBType,
			Nullable: c.Nullable,
		})
	}
	return res
}

func DiffColumn(ops []AlterTableSuboperation, c1, c2 *schema.Column) []AlterTableSuboperation {
	if c1.Nullable && !c2.Nullable {
		ops = append(ops, &AlterTableSetNotNull{Name: c1.Name})
	}
	if !c1.Nullable && c2.Nullable {
		ops = append(ops, &AlterTableSetNull{Name: c1.Name})
	}
	if c1.DBType != c2.DBType {
		ops = append(ops, &AlterTableSetType{
			Name: c1.Name,
			Type: c2.DBType,
		})
	}
	return ops
}