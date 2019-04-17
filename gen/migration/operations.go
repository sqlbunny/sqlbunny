package migration

import (
	"fmt"

	"github.com/kernelpayments/sqlbunny/runtime/migration"
	"github.com/kernelpayments/sqlbunny/schema"
)

func ApplyOperation(o migration.Operation, d *schema.Schema) {
	switch o := o.(type) {
	case migration.CreateTableOperation:
		if _, ok := d.Models[o.Name]; ok {
			panic("CreateTableOperation on already-existing table: " + o.Name)
		}
		cols := make([]*schema.Column, len(o.Columns))
		for i, c := range o.Columns {
			cols[i] = &schema.Column{
				Name:       c.Name,
				Nullable:   c.Nullable,
				SQLType:    c.Type,
				SQLDefault: c.Default,
			}
		}
		d.Models[o.Name] = &schema.Model{
			Name:    o.Name,
			Columns: cols,
		}

	case migration.DropTableOperation:
		if _, ok := d.Models[o.Name]; !ok {
			panic("DropTableOperation on non-existing table: " + o.Name)
		}
		delete(d.Models, o.Name)

	case migration.RenameTableOperation:
		m, ok := d.Models[o.OldName]
		if !ok {
			panic("RenameTableOperation on non-existing table: " + o.OldName)
		}
		if _, ok := d.Models[o.NewName]; ok {
			panic("RenameTableOperation new table name already exists: " + o.NewName)
		}

		delete(d.Models, o.OldName)
		d.Models[o.NewName] = m
		m.Name = o.NewName

		for _, m2 := range d.Models {
			for _, fk := range m2.ForeignKeys {
				if fk.ForeignModel == o.OldName {
					fk.ForeignModel = o.NewName
				}
			}
		}

	case migration.AlterTableOperation:
		t, ok := d.Models[o.Name]
		if !ok {
			panic("AlterTableOperation on non-existing table: " + o.Name)
		}
		for _, op := range o.Ops {
			ApplyAlterTableOperation(op, d, t)
		}

	case migration.CreateIndexOperation:
		t, ok := d.Models[o.Name]
		if !ok {
			panic("CreateIndexOperation on non-existing table: " + o.Name)
		}
		if t.FindIndex(o.IndexName) != nil {
			panic(fmt.Sprintf("CreateIndexOperation index already exists: table %s, index %s ", o.Name, o.IndexName))
		}
		t.Indexes = append(t.Indexes, &schema.Index{
			Name:    o.IndexName,
			Columns: o.Columns,
		})

	case migration.DropIndexOperation:
		t, ok := d.Models[o.Name]
		if !ok {
			panic("DropIndexOperation on non-existing table: " + o.Name)
		}
		if t.FindIndex(o.IndexName) == nil {
			panic(fmt.Sprintf("DropIndexOperation index doesn't exist: table %s, index %s ", o.Name, o.IndexName))
		}
		t.DeleteIndex(o.IndexName)

	case migration.RenameColumnOperation:
		m, ok := d.Models[o.Name]
		if !ok {
			panic("RenameColumnOperation on non-existing table: " + o.Name)
		}

		c := m.FindColumn(o.OldColumnName)
		if c == nil {
			panic(fmt.Sprintf("RenameColumnOperation non-existing: table %s, column %s", m.Name, o.OldColumnName))
		}

		c.Name = o.NewColumnName

		if m.PrimaryKey != nil {
			for i := range m.PrimaryKey.Columns {
				if m.PrimaryKey.Columns[i] == o.OldColumnName {
					m.PrimaryKey.Columns[i] = o.NewColumnName
				}
			}
		}
		for _, idx := range m.Indexes {
			for i := range idx.Columns {
				if idx.Columns[i] == o.OldColumnName {
					idx.Columns[i] = o.NewColumnName
				}
			}
		}
		for _, idx := range m.Uniques {
			for i := range idx.Columns {
				if idx.Columns[i] == o.OldColumnName {
					idx.Columns[i] = o.NewColumnName
				}
			}
		}

		for _, m2 := range d.Models {
			for _, fk := range m2.ForeignKeys {
				if fk.ForeignModel == m.Name && fk.ForeignColumn == o.OldColumnName {
					fk.ForeignColumn = o.NewColumnName
				}
			}
		}

	case migration.SQLOperation:
		// do nothing.

	default:
		panic(fmt.Sprintf("Unknown operation type: %T", o))
	}
}

func ApplyAlterTableOperation(o migration.AlterTableSuboperation, d *schema.Schema, m *schema.Model) {
	switch o := o.(type) {
	case migration.AlterTableAddColumn:
		if c := m.FindColumn(o.Name); c != nil {
			panic(fmt.Sprintf("AlterTableAddColumn already-existing: table %s, column %s", m.Name, o.Name))
		}
		m.Columns = append(m.Columns, &schema.Column{
			Name:       o.Name,
			SQLType:    o.Type,
			SQLDefault: o.Default,
			Nullable:   o.Nullable,
		})

	case migration.AlterTableDropColumn:
		if c := m.FindColumn(o.Name); c == nil {
			panic(fmt.Sprintf("AlterTableDropColumn non-existing: table %s, column %s", m.Name, o.Name))
		}
		m.DeleteColumn(o.Name)

	case migration.AlterTableCreatePrimaryKey:
		if m.PrimaryKey != nil {
			panic(fmt.Sprintf("AlterTableCreatePrimaryKey on a model already with primary key: %s", m.Name))
		}
		m.PrimaryKey = &schema.PrimaryKey{
			Columns: o.Columns,
		}

	case migration.AlterTableDropPrimaryKey:
		if m.PrimaryKey == nil {
			panic(fmt.Sprintf("AlterTableDropPrimaryKey on a model already without primary key: %s", m.Name))
		}
		m.PrimaryKey = nil

	case migration.AlterTableCreateUnique:
		idx := m.FindUnique(o.Name)
		if idx != nil {
			panic(fmt.Sprintf("AlterTableCreateUnique unique already exists: table %s, unique %s ", m.Name, o.Name))
		}
		m.Uniques = append(m.Uniques, &schema.Unique{
			Name:    o.Name,
			Columns: o.Columns,
		})

	case migration.AlterTableDropUnique:
		idx := m.FindUnique(o.Name)
		if idx == nil {
			panic(fmt.Sprintf("AlterTableDropUnique unique doesn't exist: table %s, unique %s ", m.Name, o.Name))
		}
		m.DeleteUnique(o.Name)

	case migration.AlterTableCreateForeignKey:
		idx := m.FindForeignKey(o.Name)
		if idx != nil {
			panic(fmt.Sprintf("AlterTableCreateForeignKey ForeignKey already exists: table %s, ForeignKey %s ", m.Name, o.Name))
		}
		if len(o.Columns) != len(o.ForeignColumns) {
			panic(fmt.Sprintf("AlterTableCreateForeignKey lengths of Columns and ForeignColumns must match: table %s, ForeignKey %s ", m.Name, o.Name))
		}
		if len(o.Columns) != 1 || len(o.ForeignColumns) != 1 {
			panic(fmt.Sprintf("AlterTableCreateForeignKey multi-column FKs are not yet supported: table %s, ForeignKey %s ", m.Name, o.Name))
		}
		m.ForeignKeys = append(m.ForeignKeys, &schema.ForeignKey{
			Name:          o.Name,
			Model:         m.Name,
			Column:        o.Columns[0],
			ForeignModel:  o.ForeignTable,
			ForeignColumn: o.ForeignColumns[0],
		})

	case migration.AlterTableDropForeignKey:
		idx := m.FindForeignKey(o.Name)
		if idx == nil {
			panic(fmt.Sprintf("AlterTableDropForeignKey ForeignKey doesn't exist: table %s, ForeignKey %s ", m.Name, o.Name))
		}
		m.DeleteForeignKey(o.Name)

	case migration.AlterTableSetNotNull:
		c := m.FindColumn(o.Name)
		if c == nil {
			panic(fmt.Sprintf("AlterTableSetNotNull column doesn't exist: table %s, column %s ", m.Name, o.Name))
		}
		c.Nullable = false

	case migration.AlterTableSetNull:
		c := m.FindColumn(o.Name)
		if c == nil {
			panic(fmt.Sprintf("AlterTableSetNull column doesn't exist: table %s, column %s ", m.Name, o.Name))
		}
		c.Nullable = true

	case migration.AlterTableSetDefault:
		c := m.FindColumn(o.Name)
		if c == nil {
			panic(fmt.Sprintf("AlterTableSetDefault column doesn't exist: table %s, column %s ", m.Name, o.Name))
		}
		c.SQLDefault = o.Default

	case migration.AlterTableDropDefault:
		c := m.FindColumn(o.Name)
		if c == nil {
			panic(fmt.Sprintf("AlterTableDropDefault column doesn't exist: table %s, column %s ", m.Name, o.Name))
		}
		c.SQLDefault = ""

	case migration.AlterTableSetType:
		c := m.FindColumn(o.Name)
		if c == nil {
			panic(fmt.Sprintf("AlterTableSetType column doesn't exist: table %s, column %s ", m.Name, o.Name))
		}
		c.SQLType = o.Type

	default:
		panic(fmt.Sprintf("Unknown alter table suboperation type: %T", o))
	}
}
