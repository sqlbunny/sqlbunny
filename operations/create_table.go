package operations

import (
	"fmt"
	"io"
	"strings"

	"github.com/sqlbunny/sqlschema/schema"
)

type Column struct {
	Name     string
	Type     string
	Default  string
	Nullable bool
}

func (c Column) Dump(w io.Writer) {
	fmt.Fprintf(w, "operations.Column{Name: %s, Type: %s, Default: %s, Nullable: %s}", esc(c.Name), esc(c.Type), esc(c.Default), dumpBool(c.Nullable))
}

type CreateTable struct {
	SchemaName string
	TableName  string
	Columns    []Column
}

func (o CreateTable) GetSQL() string {
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
	return fmt.Sprintf("CREATE TABLE %s (\n%s\n)", sqlName(o.SchemaName, o.TableName), strings.Join(x, ",\n"))
}

func (o CreateTable) Dump(w io.Writer) {
	fmt.Fprint(w, "operations.CreateTable {\n")
	fmt.Fprint(w, "SchemaName: "+esc(o.SchemaName)+",\n")
	fmt.Fprint(w, "TableName: "+esc(o.TableName)+",\n")
	fmt.Fprint(w, "Columns: []operations.Column{\n")
	for _, c := range o.Columns {
		c.Dump(w)
		fmt.Fprint(w, ",\n")
	}
	fmt.Fprint(w, "},\n")
	fmt.Fprint(w, "}")
}

func (o CreateTable) Apply(d *schema.Database) error {
	s, ok := d.Schemas[o.SchemaName]
	if !ok {
		return fmt.Errorf("no such schema: %s", o.SchemaName)
	}
	if _, ok := s.Tables[o.TableName]; ok {
		return fmt.Errorf("table already exists: %s", o.TableName)
	}

	t := schema.NewTable()
	for _, c := range o.Columns {
		t.Columns[c.Name] = &schema.Column{
			Nullable: c.Nullable,
			Type:     c.Type,
			Default:  c.Default,
		}
	}
	s.Tables[o.TableName] = t

	return nil
}
