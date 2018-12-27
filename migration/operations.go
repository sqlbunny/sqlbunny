package migration

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/kernelpayments/sqlbunny/runtime/bunny"
	"github.com/kernelpayments/sqlbunny/schema"
)

func getDefaultForType(dbType string) string {
	switch dbType {
	case "smallint", "integer", "bigint", "decimal", "numeric", "real", "double precision", "int2", "int4", "int8":
		return "0"
	case "boolean":
		return "false"
	case "varchar", "text":
		return "''"
	case "bytea":
		return "''"
	case "jsonb":
		return "'{}'"
	}

	// For arrays, the default is an empty array.
	if strings.HasSuffix(dbType, "[]") {
		return "'{}'"
	}

	return ""
}

type Operation interface {
	Run(ctx context.Context) error
	Apply(d *schema.Schema)
	Dump(buf *bytes.Buffer)
}

type OperationList []Operation

func (o OperationList) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.OperationList {\n")
	for _, op := range o {
		op.Dump(buf)
		buf.WriteString(",\n")
	}
	buf.WriteString("}")
}

type Column struct {
	Name     string
	Type     string
	Nullable bool
}

func (o Column) Dump(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("migration.Column{Name: \"%s\", Type: \"%s\", Nullable: %s}", o.Name, o.Type, dumpBool(o.Nullable)))
}

type CreateTableOperation struct {
	Name    string
	Columns []Column
}

func (o CreateTableOperation) Run(ctx context.Context) error {
	var x []string
	for _, c := range o.Columns {
		var n string
		if !c.Nullable {
			n = "NOT NULL"
			d := getDefaultForType(c.Type)
			if d != "" {
				n += " DEFAULT " + d
			}
		}
		x = append(x, fmt.Sprintf("\"%s\" %s %s", c.Name, c.Type, n))
	}
	q := fmt.Sprintf("CREATE TABLE \"%s\" (%s)", o.Name, strings.Join(x, ","))
	_, err := bunny.Exec(ctx, q)
	return err
}

func (o CreateTableOperation) Apply(d *schema.Schema) {
	if _, ok := d.Models[o.Name]; ok {
		panic("CreateTableOperation on already-existing table: " + o.Name)
	}
	cols := make([]*schema.Column, len(o.Columns))
	for i, c := range o.Columns {
		cols[i] = &schema.Column{
			Name:     c.Name,
			Nullable: c.Nullable,
			DBType:   c.Type,
		}
	}
	d.Models[o.Name] = &schema.Model{
		Name:    o.Name,
		Columns: cols,
	}
}

func (o CreateTableOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.CreateTableOperation {\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("Columns: []migration.Column{\n")
	for _, c := range o.Columns {
		c.Dump(buf)
		buf.WriteString(",\n")
	}
	buf.WriteString("},\n")
	buf.WriteString("}")
}

type DropTableOperation struct {
	Name string
}

func (o DropTableOperation) Run(ctx context.Context) error {
	q := fmt.Sprintf("DROP TABLE \"%s\"", o.Name)
	_, err := bunny.Exec(ctx, q)
	return err
}
func (o DropTableOperation) Apply(d *schema.Schema) {
	if _, ok := d.Models[o.Name]; !ok {
		panic("DropTableOperation on non-existing table: " + o.Name)
	}
	delete(d.Models, o.Name)
}

func (o DropTableOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.DropTableOperation {\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("}")
}

type AlterTableSuboperation interface {
	AlterTableSQL(ato *AlterTableOperation) string
	Apply(d *schema.Schema, m *schema.Model)
	Dump(buf *bytes.Buffer)
}

type AlterTableAddColumn struct {
	Name     string
	Type     string
	Nullable bool
}

func (o *AlterTableAddColumn) AlterTableSQL(ato *AlterTableOperation) string {
	var n string
	if !o.Nullable {
		n = "NOT NULL"
		d := getDefaultForType(o.Type)
		if d != "" {
			n += " DEFAULT " + d
		}
	}
	return fmt.Sprintf("ADD COLUMN \"%s\" %s %s", o.Name, o.Type, n)
}
func (o *AlterTableAddColumn) Apply(d *schema.Schema, m *schema.Model) {
	if c := m.FindColumn(o.Name); c != nil {
		panic(fmt.Sprintf("AlterTableAddColumn already-existing: table %s, column %s", m.Name, o.Name))
	}
	m.Columns = append(m.Columns, &schema.Column{
		Name:     o.Name,
		DBType:   o.Type,
		Nullable: o.Nullable,
	})
}
func (o *AlterTableAddColumn) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableAddColumn {\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("Type: \"" + o.Type + "\",\n")
	buf.WriteString("Nullable: " + dumpBool(o.Nullable) + ",\n")
	buf.WriteString("}")
}

func dumpBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

type AlterTableDropColumn struct {
	Name string
}

func (o *AlterTableDropColumn) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("DROP COLUMN \"%s\"", o.Name)
}
func (o *AlterTableDropColumn) Apply(d *schema.Schema, m *schema.Model) {
	if c := m.FindColumn(o.Name); c == nil {
		panic(fmt.Sprintf("AlterTableDropColumn non-existing: table %s, column %s", m.Name, o.Name))
	}
	m.DeleteColumn(o.Name)
}
func (o *AlterTableDropColumn) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableDropColumn {Name: \"" + o.Name + "\"}")
}

func columnList(columns []string) string {
	var buf bytes.Buffer
	for i, c := range columns {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("\"")
		buf.WriteString(c)
		buf.WriteString("\"")
	}
	return buf.String()
}

type AlterTableCreatePrimaryKey struct {
	Columns []string
}

func (o *AlterTableCreatePrimaryKey) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s_pkey\" PRIMARY KEY (%s)", ato.Name, columnList(o.Columns))
}
func (o *AlterTableCreatePrimaryKey) Apply(d *schema.Schema, m *schema.Model) {
	if m.PrimaryKey != nil {
		panic(fmt.Sprintf("AlterTableCreatePrimaryKey on a model already with primary key: %s", m.Name))
	}
	m.PrimaryKey = &schema.PrimaryKey{
		Columns: o.Columns,
	}
}
func (o *AlterTableCreatePrimaryKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableCreatePrimaryKey{\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("}")
}

type AlterTableDropPrimaryKey struct {
}

func (o *AlterTableDropPrimaryKey) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s_pkey\"", ato.Name)
}
func (o *AlterTableDropPrimaryKey) Apply(d *schema.Schema, m *schema.Model) {
	if m.PrimaryKey == nil {
		panic(fmt.Sprintf("AlterTableDropPrimaryKey on a model already without primary key: %s", m.Name))
	}
	m.PrimaryKey = nil
}
func (o *AlterTableDropPrimaryKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableDropPrimaryKey{}")
}

type AlterTableCreateUnique struct {
	Name    string
	Columns []string
}

func (o *AlterTableCreateUnique) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s\" UNIQUE (%s)", o.Name, columnList(o.Columns))
}
func (o *AlterTableCreateUnique) Apply(d *schema.Schema, m *schema.Model) {
	idx := m.FindUnique(o.Name)
	if idx != nil {
		panic(fmt.Sprintf("AlterTableCreateUnique unique already exists: table %s, unique %s ", m.Name, o.Name))
	}
	m.Uniques = append(m.Uniques, &schema.Unique{
		Name:    o.Name,
		Columns: o.Columns,
	})
}
func (o *AlterTableCreateUnique) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableCreateUnique{\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("}")
}

type AlterTableDropUnique struct {
	Name string
}

func (o *AlterTableDropUnique) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s\"", o.Name)
}
func (o *AlterTableDropUnique) Apply(d *schema.Schema, m *schema.Model) {
	idx := m.FindUnique(o.Name)
	if idx == nil {
		panic(fmt.Sprintf("AlterTableDropUnique unique doesn't exist: table %s, unique %s ", m.Name, o.Name))
	}
	m.DeleteUnique(o.Name)
}
func (o *AlterTableDropUnique) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableDropUnique{\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("}")
}

type AlterTableCreateForeignKey struct {
	Name           string
	Columns        []string
	ForeignTable   string
	ForeignColumns []string
}

func (o *AlterTableCreateForeignKey) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s\" FOREIGN KEY (%s) REFERENCES \"%s\" (%s)", o.Name, columnList(o.Columns), o.ForeignTable, columnList(o.ForeignColumns))
}
func (o *AlterTableCreateForeignKey) Apply(d *schema.Schema, m *schema.Model) {
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
}
func (o *AlterTableCreateForeignKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableCreateForeignKey{\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("ForeignTable: \"" + o.ForeignTable + "\",\n")
	buf.WriteString("ForeignColumns: []string{" + columnList(o.ForeignColumns) + "},\n")
	buf.WriteString("}")
}

type AlterTableDropForeignKey struct {
	Name string
}

func (o *AlterTableDropForeignKey) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s\"", o.Name)
}
func (o *AlterTableDropForeignKey) Apply(d *schema.Schema, m *schema.Model) {
	idx := m.FindForeignKey(o.Name)
	if idx == nil {
		panic(fmt.Sprintf("AlterTableDropForeignKey ForeignKey doesn't exist: table %s, ForeignKey %s ", m.Name, o.Name))
	}
	m.DeleteForeignKey(o.Name)
}
func (o *AlterTableDropForeignKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableDropForeignKey{\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("}")
}

type AlterTableSetNotNull struct {
	Name string
}

func (o *AlterTableSetNotNull) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" SET NOT NULL", o.Name)
}
func (o *AlterTableSetNotNull) Apply(d *schema.Schema, m *schema.Model) {
	c := m.FindColumn(o.Name)
	if c == nil {
		panic(fmt.Sprintf("AlterTableSetNotNull column doesn't exist: table %s, column %s ", m.Name, o.Name))
	}
	c.Nullable = false
}
func (o *AlterTableSetNotNull) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableSetNotNull{Name: \"" + o.Name + "\"}")
}

type AlterTableSetNull struct {
	Name string
}

func (o *AlterTableSetNull) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" DROP NOT NULL", o.Name)
}
func (o *AlterTableSetNull) Apply(d *schema.Schema, m *schema.Model) {
	c := m.FindColumn(o.Name)
	if c == nil {
		panic(fmt.Sprintf("AlterTableSetNull column doesn't exist: table %s, column %s ", m.Name, o.Name))
	}
	c.Nullable = true
}
func (o *AlterTableSetNull) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableSetNull{Name: \"" + o.Name + "\"}")
}

type AlterTableSetType struct {
	Name string
	Type string
}

func (o *AlterTableSetType) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" TYPE %s", o.Name, o.Type)
}
func (o *AlterTableSetType) Apply(d *schema.Schema, m *schema.Model) {
	c := m.FindColumn(o.Name)
	if c == nil {
		panic(fmt.Sprintf("AlterTableSetType column doesn't exist: table %s, column %s ", m.Name, o.Name))
	}
	c.DBType = o.Type
}
func (o *AlterTableSetType) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.AlterTableSetType{\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("Type: \"" + o.Type + "\",\n")
	buf.WriteString("}")
}

type AlterTableOperation struct {
	Name string
	Ops  []AlterTableSuboperation
}

func (o AlterTableOperation) Run(ctx context.Context) error {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("ALTER TABLE \"%s\" ", o.Name))
	first := true
	for _, op := range o.Ops {
		if !first {
			buf.WriteString(", ")
		}
		buf.WriteString(op.AlterTableSQL(&o))
		first = false
	}
	_, err := bunny.Exec(ctx, buf.String())
	return err
}
func (o AlterTableOperation) Apply(d *schema.Schema) {
	t, ok := d.Models[o.Name]
	if !ok {
		panic("AlterTableOperation on non-existing table: " + o.Name)
	}
	for _, op := range o.Ops {
		op.Apply(d, t)
	}
}

func (o AlterTableOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableOperation {\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("Ops: []migration.AlterTableSuboperation{\n")
	for _, op := range o.Ops {
		op.Dump(buf)
		buf.WriteString(",\n")
	}
	buf.WriteString("},\n")
	buf.WriteString("}")
}

type CreateIndexOperation struct {
	Name      string
	IndexName string
	Columns   []string
}

func (o CreateIndexOperation) Run(ctx context.Context) error {
	q := fmt.Sprintf("CREATE INDEX CONCURRENTLY \"%s\" ON \"%s\" (%s) ", o.IndexName, o.Name, columnList(o.Columns))
	_, err := bunny.Exec(ctx, q)
	return err
}
func (o CreateIndexOperation) Apply(d *schema.Schema) {
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
}

func (o CreateIndexOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.CreateIndexOperation {\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("IndexName: \"" + o.IndexName + "\",\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("}")
}

type RenameColumnOperation struct {
	Name          string
	OldColumnName string
	NewColumnName string
}

func (o RenameColumnOperation) Run(ctx context.Context) error {
	q := fmt.Sprintf("ALTER TABLE \"%s\" RENAME COLUMN \"%s\" TO \"%s\"", o.Name, o.OldColumnName, o.NewColumnName)
	_, err := bunny.Exec(ctx, q)
	return err
}

func (o RenameColumnOperation) Apply(d *schema.Schema) {
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
}

func (o RenameColumnOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("&migration.RenameColumnOperation {Name: \"" + o.Name + "\", OldColumnName: \"" + o.OldColumnName + "\", NewColumnName: \"" + o.NewColumnName + "\"}")
}

type DropIndexOperation struct {
	Name      string
	IndexName string
}

func (o DropIndexOperation) Run(ctx context.Context) error {
	q := fmt.Sprintf("DROP INDEX \"%s\"", o.IndexName)
	_, err := bunny.Exec(ctx, q)
	return err
}
func (o DropIndexOperation) Apply(d *schema.Schema) {
	t, ok := d.Models[o.Name]
	if !ok {
		panic("DropIndexOperation on non-existing table: " + o.Name)
	}
	if t.FindIndex(o.IndexName) == nil {
		panic(fmt.Sprintf("DropIndexOperation index doesn't exist: table %s, index %s ", o.Name, o.IndexName))
	}
	t.DeleteIndex(o.IndexName)
}

func (o DropIndexOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.DropIndexOperation {\n")
	buf.WriteString("Name: \"" + o.Name + "\",\n")
	buf.WriteString("IndexName: \"" + o.IndexName + "\",\n")
	buf.WriteString("}")
}

type SQLOperation struct {
	SQL string
}

func (o SQLOperation) Run(ctx context.Context) error {
	_, err := bunny.Exec(ctx, o.SQL)
	return err
}
func (o SQLOperation) Apply(d *schema.Schema) {
}

func (o SQLOperation) Dump(buf *bytes.Buffer) {
	panic("SQLOperation can't Dump")
}
