package migration

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/kernelpayments/sqlbunny/runtime/bunny"
)

func esc(s string) string {
	return fmt.Sprintf("%#v", s)
}

type Operation interface {
	Run(ctx context.Context) error
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
	Default  string
	Nullable bool
}

func (o Column) Dump(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf("migration.Column{Name: %s, Type: %s, Default: %s, Nullable: %s}", esc(o.Name), esc(o.Type), esc(o.Default), dumpBool(o.Nullable)))
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
			n = " NOT NULL"
		}
		var d string
		if c.Default != "" {
			d = " DEFAULT " + c.Default
		}
		x = append(x, fmt.Sprintf("    \"%s\" %s%s%s", c.Name, c.Type, n, d))
	}
	q := fmt.Sprintf("CREATE TABLE \"%s\" (\n%s\n)", o.Name, strings.Join(x, ",\n"))
	_, err := bunny.Exec(ctx, q)
	return err
}

func (o CreateTableOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.CreateTableOperation {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
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

func (o DropTableOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.DropTableOperation {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("}")
}

type RenameTableOperation struct {
	OldName string
	NewName string
}

func (o RenameTableOperation) Run(ctx context.Context) error {
	q := fmt.Sprintf("ALTER TABLE \"%s\" RENAME TO \"%s\"", o.OldName, o.NewName)
	_, err := bunny.Exec(ctx, q)
	return err
}

func (o RenameTableOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.RenameTableOperation {\n")
	buf.WriteString("OldName: " + esc(o.OldName) + ",\n")
	buf.WriteString("NewName: " + esc(o.NewName) + ",\n")
	buf.WriteString("}")
}

type AlterTableSuboperation interface {
	AlterTableSQL(ato *AlterTableOperation) string
	Dump(buf *bytes.Buffer)
}

type AlterTableAddColumn struct {
	Name     string
	Type     string
	Default  string
	Nullable bool
}

func (o AlterTableAddColumn) AlterTableSQL(ato *AlterTableOperation) string {
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
	buf.WriteString("migration.AlterTableAddColumn {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Type: " + esc(o.Type) + ",\n")
	buf.WriteString("Default: " + esc(o.Default) + ",\n")
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

func (o AlterTableDropColumn) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("DROP COLUMN \"%s\"", o.Name)
}

func (o AlterTableDropColumn) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableDropColumn {Name: " + esc(o.Name) + "}")
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

func (o AlterTableCreatePrimaryKey) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s_pkey\" PRIMARY KEY (%s)", ato.Name, columnList(o.Columns))
}

func (o AlterTableCreatePrimaryKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableCreatePrimaryKey{\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("}")
}

type AlterTableDropPrimaryKey struct {
}

func (o AlterTableDropPrimaryKey) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s_pkey\"", ato.Name)
}

func (o AlterTableDropPrimaryKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableDropPrimaryKey{}")
}

type AlterTableCreateUnique struct {
	Name    string
	Columns []string
}

func (o AlterTableCreateUnique) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s\" UNIQUE (%s)", o.Name, columnList(o.Columns))
}

func (o AlterTableCreateUnique) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableCreateUnique{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("}")
}

type AlterTableDropUnique struct {
	Name string
}

func (o AlterTableDropUnique) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s\"", o.Name)
}

func (o AlterTableDropUnique) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableDropUnique{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("}")
}

type AlterTableCreateForeignKey struct {
	Name           string
	Columns        []string
	ForeignTable   string
	ForeignColumns []string
}

func (o AlterTableCreateForeignKey) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ADD CONSTRAINT \"%s\" FOREIGN KEY (%s) REFERENCES \"%s\" (%s)", o.Name, columnList(o.Columns), o.ForeignTable, columnList(o.ForeignColumns))
}

func (o AlterTableCreateForeignKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableCreateForeignKey{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Columns: []string{" + columnList(o.Columns) + "},\n")
	buf.WriteString("ForeignTable: " + esc(o.ForeignTable) + ",\n")
	buf.WriteString("ForeignColumns: []string{" + columnList(o.ForeignColumns) + "},\n")
	buf.WriteString("}")
}

type AlterTableDropForeignKey struct {
	Name string
}

func (o AlterTableDropForeignKey) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("DROP CONSTRAINT \"%s\"", o.Name)
}

func (o AlterTableDropForeignKey) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableDropForeignKey{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("}")
}

type AlterTableSetNotNull struct {
	Name string
}

func (o AlterTableSetNotNull) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" SET NOT NULL", o.Name)
}

func (o AlterTableSetNotNull) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableSetNotNull{Name: " + esc(o.Name) + "}")
}

type AlterTableSetNull struct {
	Name string
}

func (o AlterTableSetNull) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" DROP NOT NULL", o.Name)
}

func (o AlterTableSetNull) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableSetNull{Name: " + esc(o.Name) + "}")
}

type AlterTableSetDefault struct {
	Name    string
	Default string
}

func (o AlterTableSetDefault) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" SET DEFAULT %s", o.Name, o.Default)
}

func (o AlterTableSetDefault) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableSetDefault{Name: " + esc(o.Name) + ", Default: " + esc(o.Default) + "}")
}

type AlterTableDropDefault struct {
	Name string
}

func (o AlterTableDropDefault) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" DROP DEFAULT", o.Name)
}

func (o AlterTableDropDefault) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableDropDefault{Name: " + esc(o.Name) + "}")
}

type AlterTableSetType struct {
	Name string
	Type string
}

func (o AlterTableSetType) AlterTableSQL(ato *AlterTableOperation) string {
	return fmt.Sprintf("ALTER COLUMN \"%s\" TYPE %s", o.Name, o.Type)
}

func (o AlterTableSetType) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableSetType{\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("Type: " + esc(o.Type) + ",\n")
	buf.WriteString("}")
}

type AlterTableOperation struct {
	Name string
	Ops  []AlterTableSuboperation
}

func (o AlterTableOperation) Run(ctx context.Context) error {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("ALTER TABLE \"%s\"\n", o.Name))
	first := true
	for _, op := range o.Ops {
		if !first {
			buf.WriteString(",\n")
		}
		buf.WriteString("    ")
		buf.WriteString(op.AlterTableSQL(&o))
		first = false
	}
	_, err := bunny.Exec(ctx, buf.String())
	return err
}

func (o AlterTableOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.AlterTableOperation {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
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
	q := fmt.Sprintf("CREATE INDEX CONCURRENTLY \"%s\" ON \"%s\" (%s)", o.IndexName, o.Name, columnList(o.Columns))
	_, err := bunny.Exec(ctx, q)
	return err
}

func (o CreateIndexOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.CreateIndexOperation {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("IndexName: " + esc(o.IndexName) + ",\n")
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

func (o RenameColumnOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.RenameColumnOperation {Name: " + esc(o.Name) + ", OldColumnName: " + esc(o.OldColumnName) + ", NewColumnName: " + esc(o.NewColumnName) + "}")
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

func (o DropIndexOperation) Dump(buf *bytes.Buffer) {
	buf.WriteString("migration.DropIndexOperation {\n")
	buf.WriteString("Name: " + esc(o.Name) + ",\n")
	buf.WriteString("IndexName: " + esc(o.IndexName) + ",\n")
	buf.WriteString("}")
}

type SQLOperation struct {
	SQL string
}

func (o SQLOperation) Run(ctx context.Context) error {
	_, err := bunny.Exec(ctx, o.SQL)
	return err
}

func (o SQLOperation) Dump(buf *bytes.Buffer) {
	panic("SQLOperation can't Dump")
}
