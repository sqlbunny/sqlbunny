package migration

import (
	"github.com/pkg/errors"
	"github.com/sqlbunny/sqlbunny/runtime/migration"
	"github.com/sqlbunny/sqlbunny/schema"
)

func ApplyMigration(m *migration.Migration, d *schema.Schema) error {
	for _, o := range m.Operations {
		if err := ApplyOperation(o, d); err != nil {
			return err
		}
	}
	return nil
}

func ApplyOperation(o migration.Operation, d *schema.Schema) error {
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
			return errors.Errorf("CreateIndexOperation index already exists: table %s, index %s ", o.Name, o.IndexName)
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
			return errors.Errorf("DropIndexOperation index doesn't exist: table %s, index %s ", o.Name, o.IndexName)
		}
		t.DeleteIndex(o.IndexName)

	case migration.RenameColumnOperation:
		m, ok := d.Models[o.Name]
		if !ok {
			panic("RenameColumnOperation on non-existing table: " + o.Name)
		}

		c := m.FindColumn(o.OldColumnName)
		if c == nil {
			return errors.Errorf("RenameColumnOperation non-existing: table %s, column %s", m.Name, o.OldColumnName)
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

		for _, fk := range m.ForeignKeys {
			for i := range fk.Columns {
				if fk.Columns[i] == o.OldColumnName {
					fk.Columns[i] = o.NewColumnName
				}
			}
		}

		for _, m2 := range d.Models {
			for _, fk := range m2.ForeignKeys {
				if fk.ForeignModel == m.Name {
					for _, fk := range m2.ForeignKeys {
						for i := range fk.ForeignColumns {
							if fk.ForeignColumns[i] == o.OldColumnName {
								fk.ForeignColumns[i] = o.NewColumnName
							}
						}
					}
				}
			}
		}

	case migration.SQLOperation:
		// do nothing.

	default:
		return errors.Errorf("Unknown operation type: %T", o)
	}
	return nil
}

func ApplyAlterTableOperation(o migration.AlterTableSuboperation, d *schema.Schema, m *schema.Model) error {
	switch o := o.(type) {
	case migration.AlterTableAddColumn:
		if c := m.FindColumn(o.Name); c != nil {
			return errors.Errorf("AlterTableAddColumn already-existing: table %s, column %s", m.Name, o.Name)
		}
		m.Columns = append(m.Columns, &schema.Column{
			Name:       o.Name,
			SQLType:    o.Type,
			SQLDefault: o.Default,
			Nullable:   o.Nullable,
		})

	case migration.AlterTableDropColumn:
		if c := m.FindColumn(o.Name); c == nil {
			return errors.Errorf("AlterTableDropColumn non-existing: table %s, column %s", m.Name, o.Name)
		}
		m.DeleteColumn(o.Name)

	case migration.AlterTableCreatePrimaryKey:
		if m.PrimaryKey != nil {
			return errors.Errorf("AlterTableCreatePrimaryKey on a model already with primary key: %s", m.Name)
		}
		m.PrimaryKey = &schema.PrimaryKey{
			Columns: o.Columns,
		}

	case migration.AlterTableDropPrimaryKey:
		if m.PrimaryKey == nil {
			return errors.Errorf("AlterTableDropPrimaryKey on a model already without primary key: %s", m.Name)
		}
		m.PrimaryKey = nil

	case migration.AlterTableCreateUnique:
		idx := m.FindUnique(o.Name)
		if idx != nil {
			return errors.Errorf("AlterTableCreateUnique unique already exists: table %s, unique %s ", m.Name, o.Name)
		}
		m.Uniques = append(m.Uniques, &schema.Unique{
			Name:    o.Name,
			Columns: o.Columns,
		})

	case migration.AlterTableDropUnique:
		idx := m.FindUnique(o.Name)
		if idx == nil {
			return errors.Errorf("AlterTableDropUnique unique doesn't exist: table %s, unique %s ", m.Name, o.Name)
		}
		m.DeleteUnique(o.Name)

	case migration.AlterTableCreateForeignKey:
		idx := m.FindForeignKey(o.Name)
		if idx != nil {
			return errors.Errorf("AlterTableCreateForeignKey ForeignKey already exists: table %s, ForeignKey %s ", m.Name, o.Name)
		}
		if len(o.Columns) != len(o.ForeignColumns) {
			return errors.Errorf("AlterTableCreateForeignKey lengths of Columns and ForeignColumns must match: table %s, ForeignKey %s ", m.Name, o.Name)
		}
		m.ForeignKeys = append(m.ForeignKeys, &schema.ForeignKey{
			Name:           o.Name,
			Model:          m.Name,
			Columns:        o.Columns,
			ForeignModel:   o.ForeignTable,
			ForeignColumns: o.ForeignColumns,
		})

	case migration.AlterTableDropForeignKey:
		idx := m.FindForeignKey(o.Name)
		if idx == nil {
			return errors.Errorf("AlterTableDropForeignKey ForeignKey doesn't exist: table %s, ForeignKey %s ", m.Name, o.Name)
		}
		m.DeleteForeignKey(o.Name)

	case migration.AlterTableSetNotNull:
		c := m.FindColumn(o.Name)
		if c == nil {
			return errors.Errorf("AlterTableSetNotNull column doesn't exist: table %s, column %s ", m.Name, o.Name)
		}
		c.Nullable = false

	case migration.AlterTableSetNull:
		c := m.FindColumn(o.Name)
		if c == nil {
			return errors.Errorf("AlterTableSetNull column doesn't exist: table %s, column %s ", m.Name, o.Name)
		}
		c.Nullable = true

	case migration.AlterTableSetDefault:
		c := m.FindColumn(o.Name)
		if c == nil {
			return errors.Errorf("AlterTableSetDefault column doesn't exist: table %s, column %s ", m.Name, o.Name)
		}
		c.SQLDefault = o.Default

	case migration.AlterTableDropDefault:
		c := m.FindColumn(o.Name)
		if c == nil {
			return errors.Errorf("AlterTableDropDefault column doesn't exist: table %s, column %s ", m.Name, o.Name)
		}
		c.SQLDefault = ""

	case migration.AlterTableSetType:
		c := m.FindColumn(o.Name)
		if c == nil {
			return errors.Errorf("AlterTableSetType column doesn't exist: table %s, column %s ", m.Name, o.Name)
		}
		c.SQLType = o.Type

	default:
		return errors.Errorf("Unknown alter table suboperation type: %T", o)
	}

	return nil
}
