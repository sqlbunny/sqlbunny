package operations

import (
	"fmt"
	"io"

	"github.com/sqlbunny/sqlschema/schema"
)

type AlterTableSuboperation interface {
	GetAlterTableSQL(ato *AlterTable) string
	Dump(w io.Writer)
	Apply(d *schema.Database, t *schema.Table, ato AlterTable) error
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

func (o AlterTableAddColumn) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableAddColumn {\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "Type: "+esc(o.Type)+",\n")
	fmt.Fprint(w, "Default: "+esc(o.Default)+",\n")
	fmt.Fprint(w, "Nullable: "+dumpBool(o.Nullable)+",\n")
	fmt.Fprint(w, "}")
}

func (o AlterTableAddColumn) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	if _, ok := t.Columns[o.Name]; ok {
		return fmt.Errorf("column already exists: %s", o.Name)
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

func (o AlterTableDropColumn) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableDropColumn {Name: "+esc(o.Name)+"}")
}

func (o AlterTableDropColumn) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	if _, ok := t.Columns[o.Name]; !ok {
		return fmt.Errorf("no such column: %s", o.Name)
	}
	delete(t.Columns, o.Name)
	return nil
}

type AlterTableCreatePrimaryKey struct {
	Columns []string
}

func (o AlterTableCreatePrimaryKey) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s_pkey\" PRIMARY KEY (%s)", ato.TableName, columnList(o.Columns))
}

func (o AlterTableCreatePrimaryKey) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableCreatePrimaryKey{\n")
	fmt.Fprint(w, "Columns: []string{"+columnList(o.Columns)+"},\n")
	fmt.Fprint(w, "}")
}

func (o AlterTableCreatePrimaryKey) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	if t.PrimaryKey != nil {
		return fmt.Errorf("table already has a primary key")
	}
	t.PrimaryKey = &schema.PrimaryKey{
		Columns: o.Columns,
	}
	return nil
}

type AlterTableDropPrimaryKey struct {
}

func (o AlterTableDropPrimaryKey) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s_pkey\"", ato.TableName)
}

func (o AlterTableDropPrimaryKey) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableDropPrimaryKey{}")
}

func (o AlterTableDropPrimaryKey) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	if t.PrimaryKey == nil {
		return fmt.Errorf("table does not have a primary key")
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

func (o AlterTableCreateUnique) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableCreateUnique{\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "Columns: []string{"+columnList(o.Columns)+"},\n")
	fmt.Fprint(w, "}")
}

func (o AlterTableCreateUnique) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	if _, ok := t.Uniques[o.Name]; ok {
		return fmt.Errorf("unique already exists: %s ", o.Name)
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

func (o AlterTableDropUnique) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableDropUnique{\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "}")
}

func (o AlterTableDropUnique) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	if _, ok := t.Uniques[o.Name]; !ok {
		return fmt.Errorf("no such unique: %s ", o.Name)
	}
	delete(t.Uniques, o.Name)
	return nil
}

type AlterTableCreateForeignKey struct {
	Name           string
	Columns        []string
	ForeignSchema  string
	ForeignTable   string
	ForeignColumns []string
}

func (o AlterTableCreateForeignKey) GetAlterTableSQL(ato *AlterTable) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s\" FOREIGN KEY (%s) REFERENCES %s (%s)", o.Name, columnList(o.Columns), sqlName(o.ForeignSchema, o.ForeignTable), columnList(o.ForeignColumns))
}

func (o AlterTableCreateForeignKey) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableCreateForeignKey{\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "Columns: []string{"+columnList(o.Columns)+"},\n")
	fmt.Fprint(w, "ForeignSchema: "+esc(o.ForeignSchema)+",\n")
	fmt.Fprint(w, "ForeignTable: "+esc(o.ForeignTable)+",\n")
	fmt.Fprint(w, "ForeignColumns: []string{"+columnList(o.ForeignColumns)+"},\n")
	fmt.Fprint(w, "}")
}

func (o AlterTableCreateForeignKey) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	if _, ok := t.ForeignKeys[o.Name]; ok {
		return fmt.Errorf("foreign key already exists: %s ", o.Name)
	}
	if len(o.Columns) != len(o.ForeignColumns) {
		return fmt.Errorf("lengths of Columns and ForeignColumns don't match: %s ", o.Name)
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

func (o AlterTableDropForeignKey) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableDropForeignKey{\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "}")
}

func (o AlterTableDropForeignKey) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	if _, ok := t.ForeignKeys[o.Name]; !ok {
		return fmt.Errorf("no such foreign key: %s", o.Name)
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

func (o AlterTableSetNotNull) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableSetNotNull{Name: "+esc(o.Name)+"}")
}

func (o AlterTableSetNotNull) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("no such column: %s ", o.Name)
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

func (o AlterTableSetNull) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableSetNull{Name: "+esc(o.Name)+"}")
}

func (o AlterTableSetNull) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("no such column: %s ", o.Name)
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

func (o AlterTableSetDefault) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableSetDefault{Name: "+esc(o.Name)+", Default: "+esc(o.Default)+"}")
}

func (o AlterTableSetDefault) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("no such column: %s ", o.Name)
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

func (o AlterTableDropDefault) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableDropDefault{Name: "+esc(o.Name)+"}")
}

func (o AlterTableDropDefault) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("no such column: %s ", o.Name)
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

func (o AlterTableSetType) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.AlterTableSetType{\n")
	fmt.Fprint(w, "Name: "+esc(o.Name)+",\n")
	fmt.Fprint(w, "Type: "+esc(o.Type)+",\n")
	fmt.Fprint(w, "}")
}

func (o AlterTableSetType) Apply(d *schema.Database, t *schema.Table, ato AlterTable) error {
	c, ok := t.Columns[o.Name]
	if !ok {
		return fmt.Errorf("no such column: %s ", o.Name)
	}
	c.Type = o.Type
	return nil
}
