package operations

import (
	"bytes"
	"fmt"

	"github.com/sqlbunny/sqlschema/schema"
)

type AlterTableSuboperation interface {
	GetAlterTableSQL(ato *AlterTable) string
	Dump(buf *bytes.Buffer)
	Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error
}

type AlterTableAddColumn struct {
	Name     string
	Type     string
	Default  string
	Nullable bool
}

func (o AlterTableAddColumn) GetAlterTableSQL(ato *AlterTable) string {
	var n string
	if !o.Nullable {
		n = " NOT NULL"
	}
	var d string
	if o.Default != "" {
		d = " DEFAULT " + o.Default

	}
	return fmt.Sprintf("ADD COLUMN \"%s\" %s%s%s", o.Name, o.Type, n, d)
}

func (o AlterTableAddColumn) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableAddColumn {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Type: " + esc(o.Type) + ",\n")
	buf.WriteString("Default: " + esc(o.Default) + ",\n")
	buf.WriteString("Nullable: " + dumpBool(o.Nullable) + ",\n")
	buf.WriteString("}")
}

func (o AlterTableAddColumn) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	if _, ok := t.Columns[o.Name]; ok {
		return fmt.Errorf("AlterTableAddColumn already-existing: column %s", o.Name)
	}
	t.Columns[o.Name] = &schema.Column{
		Type:     o.Type,
		Default:  o.Default,
		Nullable: o.Nullable,
	}
	return nil
}

type AlterTableDropColumn struct {
	Name string
}

func (o AlterTableDropColumn) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("DROP COLUMN \"%s\"", o.Name)
}

func (o AlterTableDropColumn) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableDropColumn {Name: " + esc(o.Name) + "}")
}

func (o AlterTableDropColumn) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	if _, ok := t.Columns[o.Name]; !ok {
		return fmt.Errorf("AlterTableDropColumn non-existing: column %s", o.Name)
	}
	delete(t.Columns, o.Name)
	return nil
}

type AlterTableCreatePrimaryKey struct {
	Columns []string
}

func (o AlterTableCreatePrimaryKey) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s_pkey\" PRIMARY KEY (%s)", ato.Name, columnList(o.Columns))
}

func (o AlterTableCreatePrimaryKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableCreatePrimaryKey{\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("}")
}

func (o AlterTableCreatePrimaryKey) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	if t.PrimaryKey != nil {
		return fmt.Errorf("AlterTableCreatePrimaryKey on a model already with primary key: %s", ato.Name)
	}
	t.PrimaryKey = &schema.PrimaryKey{
		Columns: o.Columns,
	}
	return nil
}

type AlterTableDropPrimaryKey struct {
}

func (o AlterTableDropPrimaryKey) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s_pkey\"", ato.Name)
}

func (o AlterTableDropPrimaryKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableDropPrimaryKey{}")
}

func (o AlterTableDropPrimaryKey) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	if t.PrimaryKey == nil {
		return fmt.Errorf("AlterTableDropPrimaryKey on a model already without primary key: %s", ato.Name)
	}
	t.PrimaryKey = nil
	return nil
}

type AlterTableCreateUnique struct {
	Name    string
	Columns []string
}

func (o AlterTableCreateUnique) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s\" UNIQUE (%s)", o.Name, columnList(o.Columns))
}

func (o AlterTableCreateUnique) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableCreateUnique{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("}")
}

func (o AlterTableCreateUnique) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	if _, ok := t.Uniques[o.Name]; ok {
		return fmt.Errorf("AlterTableCreateUnique unique already exists: unique %s ", o.Name)
	}
	t.Uniques[o.Name] = &schema.Unique{
		Columns: o.Columns,
	}
	return nil
}

type AlterTableDropUnique struct {
	Name string
}

func (o AlterTableDropUnique) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s\"", o.Name)
}

func (o AlterTableDropUnique) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableDropUnique{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("}")
}

func (o AlterTableDropUnique) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	if _, ok := t.Uniques[o.Name]; !ok {
		return fmt.Errorf("AlterTableDropUnique unique doesn't exist: unique %s ", o.Name)
	}
	delete(t.Uniques, o.Name)
	return nil
}

type AlterTableCreateForeignKey struct {
	Name           string
	Columns        []string
	ForeignTable   string
	ForeignColumns []string
}

func (o AlterTableCreateForeignKey) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s\" FOREIGN KEY (%s) REFERENCES \"%s\" (%s)", o.Name, columnList(o.Columns), o.ForeignTable, columnList(o.ForeignColumns))
}

func (o AlterTableCreateForeignKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableCreateForeignKey{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("ForeignTable: " + esc(o.ForeignTable) + ",\n")
	buf.WriteString("ForeignColumns: []string{" + columnList(o.ForeignColumns) + "},\n")
	buf.WriteString("}")
}

func (o AlterTableCreateForeignKey) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	if _, ok := t.ForeignKeys[o.Name]; ok {
		return fmt.Errorf("AlterTableCreateForeignKey ForeignKey already exists: ForeignKey %s ", o.Name)
	}
	if len(o.Columns) != len(o.ForeignColumns) {
		return fmt.Errorf("AlterTableCreateForeignKey lengths of Columns and ForeignColumns must match: ForeignKey %s ", o.Name)
	}
	t.ForeignKeys[o.Name] = &schema.ForeignKey{
		LocalColumns:   o.Columns,
		ForeignTable:   o.ForeignTable,
		ForeignColumns: o.ForeignColumns,
	}
	return nil
}

type AlterTableDropForeignKey struct {
	Name string
}

func (o AlterTableDropForeignKey) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s\"", o.Name)
}

func (o AlterTableDropForeignKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableDropForeignKey{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("}")
}

func (o AlterTableDropForeignKey) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	if _, ok := t.ForeignKeys[o.Name]; !ok {
		return fmt.Errorf("AlterTableDropForeignKey ForeignKey doesn't exist: ForeignKey %s ", o.Name)
	}
	delete(t.ForeignKeys, o.Name)
	return nil
}

type AlterTableSetNotNull struct {
	Name string
}

func (o AlterTableSetNotNull) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" SET NOT NULL", o.Name)
}

func (o AlterTableSetNotNull) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableSetNotNull{Name: " + esc(o.Name) + "}")
}

func (o AlterTableSetNotNull) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("AlterTableSetNotNull column doesn't exist: column %s ", o.Name)
	}
	c.Nullable = false
	return nil
}

type AlterTableSetNull struct {
	Name string
}

func (o AlterTableSetNull) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" DROP NOT NULL", o.Name)
}

func (o AlterTableSetNull) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableSetNull{Name: " + esc(o.Name) + "}")
}

func (o AlterTableSetNull) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("AlterTableSetNull column doesn't exist: column %s ", o.Name)
	}
	c.Nullable = true
	return nil
}

type AlterTableSetDefault struct {
	Name    string
	Default string
}

func (o AlterTableSetDefault) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" SET DEFAULT %s", o.Name, o.Default)
}

func (o AlterTableSetDefault) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableSetDefault{Name: " + esc(o.Name) + ", Default: " + esc(o.Default) + "}")
}

func (o AlterTableSetDefault) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("AlterTableSetDefault column doesn't exist: column %s ", o.Name)
	}
	c.Default = o.Default
	return nil
}

type AlterTableDropDefault struct {
	Name string
}

func (o AlterTableDropDefault) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" DROP DEFAULT", o.Name)
}

func (o AlterTableDropDefault) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableDropDefault{Name: " + esc(o.Name) + "}")
}

func (o AlterTableDropDefault) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("AlterTableDropDefault column doesn't exist: column %s ", o.Name)
	}
	c.Default = ""
	return nil
}

type AlterTableSetType struct {
	Name string
	Type string
}

func (o AlterTableSetType) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" TYPE %s", o.Name, o.Type)
}

func (o AlterTableSetType) Dump(buf *bytes.Buffer) {
	buf.WriteString("operations.AlterTableSetType{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Type: " + esc(o.Type) + ",\n")
	buf.WriteString("}")
}

func (o AlterTableSetType) Apply(s *schema.Schema, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("AlterTableSetType column doesn't exist: column %s ", o.Name)
	}
	c.Type = o.Type
	return nil
}
